// Package sam provides SAM resource transformers.
package sam

import (
	"fmt"
	"strings"
)

// AWS CloudFormation type for nested stacks
const (
	TypeCloudFormationStack = "AWS::CloudFormation::Stack"
)

// Application represents an AWS::Serverless::Application resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-application.html
type Application struct {
	// Location specifies the application location.
	// Can be a string (ApplicationId for SAR) or an object with ApplicationId and SemanticVersion.
	Location interface{} `json:"Location" yaml:"Location"`

	// Parameters specifies parameters to pass to the nested application.
	Parameters map[string]interface{} `json:"Parameters,omitempty" yaml:"Parameters,omitempty"`

	// NotificationArns is a list of SNS topic ARNs for stack notifications.
	NotificationArns []interface{} `json:"NotificationArns,omitempty" yaml:"NotificationArns,omitempty"`

	// Tags is a map of key-value pairs to apply to the nested stack.
	Tags map[string]string `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// TimeoutInMinutes specifies the timeout for stack creation.
	TimeoutInMinutes int `json:"TimeoutInMinutes,omitempty" yaml:"TimeoutInMinutes,omitempty"`

	// Condition is a CloudFormation condition name.
	Condition string `json:"Condition,omitempty" yaml:"Condition,omitempty"`

	// DependsOn specifies resource dependencies.
	DependsOn interface{} `json:"DependsOn,omitempty" yaml:"DependsOn,omitempty"`

	// Metadata is custom metadata for the resource.
	Metadata map[string]interface{} `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
}

// ApplicationTransformer transforms AWS::Serverless::Application to CloudFormation.
type ApplicationTransformer struct{}

// NewApplicationTransformer creates a new ApplicationTransformer.
func NewApplicationTransformer() *ApplicationTransformer {
	return &ApplicationTransformer{}
}

// Transform converts a SAM Application to CloudFormation resources.
// Returns a map of logical ID to CloudFormation resource.
func (t *ApplicationTransformer) Transform(logicalID string, app *Application, ctx *TransformContext) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Build the Stack properties
	stackProps, err := t.buildStackProperties(app, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build stack properties: %w", err)
	}

	// Build the Stack resource
	stackResource := map[string]interface{}{
		"Type":       TypeCloudFormationStack,
		"Properties": stackProps,
	}

	// Add optional resource attributes
	if app.Condition != "" {
		stackResource["Condition"] = app.Condition
	}
	if app.DependsOn != nil {
		stackResource["DependsOn"] = app.DependsOn
	}
	if app.Metadata != nil {
		stackResource["Metadata"] = app.Metadata
	}

	resources[logicalID] = stackResource

	return resources, nil
}

// buildStackProperties builds the CloudFormation Stack properties.
func (t *ApplicationTransformer) buildStackProperties(app *Application, ctx *TransformContext) (map[string]interface{}, error) {
	props := make(map[string]interface{})

	// Process Location to get TemplateURL
	templateURL, err := t.resolveTemplateURL(app.Location, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template URL: %w", err)
	}
	props["TemplateURL"] = templateURL

	// Process Parameters
	if len(app.Parameters) > 0 {
		props["Parameters"] = app.Parameters
	}

	// Process NotificationArns
	if len(app.NotificationArns) > 0 {
		props["NotificationARNs"] = app.NotificationArns
	}

	// Process Tags
	if len(app.Tags) > 0 {
		tags := make([]interface{}, 0, len(app.Tags))
		for k, v := range app.Tags {
			tags = append(tags, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		props["Tags"] = tags
	}

	// Process TimeoutInMinutes
	if app.TimeoutInMinutes > 0 {
		props["TimeoutInMinutes"] = app.TimeoutInMinutes
	}

	return props, nil
}

// resolveTemplateURL resolves the Location to a TemplateURL.
// For SAR applications, this constructs the appropriate intrinsic function.
func (t *ApplicationTransformer) resolveTemplateURL(location interface{}, ctx *TransformContext) (interface{}, error) {
	if location == nil {
		return nil, fmt.Errorf("location is required")
	}

	switch loc := location.(type) {
	case string:
		// If it's a plain string, it could be:
		// 1. An S3 URL (s3://bucket/key)
		// 2. An HTTP(S) URL
		// 3. A SAR Application ID (arn:aws:serverlessrepo:...)
		if isSARApplicationID(loc) {
			return t.buildSARTemplateURL(loc, "", ctx)
		}
		// Assume it's a direct template URL
		return loc, nil

	case map[string]interface{}:
		// Object with ApplicationId and optionally SemanticVersion
		appID, hasAppID := loc["ApplicationId"]
		semVer, hasSemVer := loc["SemanticVersion"]

		if hasAppID {
			appIDStr, ok := appID.(string)
			if !ok {
				// Could be an intrinsic function
				return t.buildSARTemplateURLWithIntrinsics(appID, semVer, ctx)
			}
			semVerStr := ""
			if hasSemVer {
				if sv, ok := semVer.(string); ok {
					semVerStr = sv
				}
			}
			return t.buildSARTemplateURL(appIDStr, semVerStr, ctx)
		}

		// Check for Bucket/Key format (S3 location)
		if bucket, hasBucket := loc["Bucket"]; hasBucket {
			key, hasKey := loc["Key"]
			if !hasKey {
				return nil, fmt.Errorf("location with Bucket requires Key")
			}
			version := loc["Version"]
			return t.buildS3TemplateURL(bucket, key, version)
		}

		return nil, fmt.Errorf("invalid Location object format")

	case map[interface{}]interface{}:
		// Convert YAML map to string map and recurse
		converted := make(map[string]interface{})
		for k, v := range loc {
			if ks, ok := k.(string); ok {
				converted[ks] = v
			}
		}
		return t.resolveTemplateURL(converted, ctx)

	default:
		return nil, fmt.Errorf("invalid Location type: %T", location)
	}
}

// isSARApplicationID checks if the string is a SAR Application ARN.
func isSARApplicationID(s string) bool {
	// SAR ARN format: arn:aws:serverlessrepo:region:account-id:applications/application-name
	return len(s) > 4 && s[0:4] == "arn:" && strings.Contains(s, ":serverlessrepo:")
}

// buildSARTemplateURL creates the TemplateURL for a SAR application.
// This uses AWS::Serverless transform macro which resolves SAR applications.
func (t *ApplicationTransformer) buildSARTemplateURL(applicationID, semanticVersion string, ctx *TransformContext) (interface{}, error) {
	// For SAR applications, we need to use the AWS::Include transform
	// or generate a reference that will be resolved by CloudFormation.
	// In practice, SAM CLI handles this during packaging.

	// Build the application ARN reference
	params := map[string]interface{}{
		"ApplicationId": applicationID,
	}
	if semanticVersion != "" {
		params["SemanticVersion"] = semanticVersion
	}

	// Return an intrinsic function that references the SAR application
	// CloudFormation will resolve this during deployment
	return map[string]interface{}{
		"Fn::Transform": map[string]interface{}{
			"Name":       "AWS::Serverless-2016-10-31",
			"Parameters": params,
		},
	}, nil
}

// buildSARTemplateURLWithIntrinsics handles SAR application IDs that use intrinsic functions.
func (t *ApplicationTransformer) buildSARTemplateURLWithIntrinsics(applicationID, semanticVersion interface{}, ctx *TransformContext) (interface{}, error) {
	params := map[string]interface{}{
		"ApplicationId": applicationID,
	}
	if semanticVersion != nil {
		params["SemanticVersion"] = semanticVersion
	}

	return map[string]interface{}{
		"Fn::Transform": map[string]interface{}{
			"Name":       "AWS::Serverless-2016-10-31",
			"Parameters": params,
		},
	}, nil
}

// buildS3TemplateURL creates a TemplateURL from S3 bucket and key.
func (t *ApplicationTransformer) buildS3TemplateURL(bucket, key, version interface{}) (interface{}, error) {
	// If all values are strings, we can construct a simple URL
	if bucketStr, ok := bucket.(string); ok {
		if keyStr, ok := key.(string); ok {
			url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketStr, keyStr)
			if version != nil {
				if versionStr, ok := version.(string); ok {
					url += "?versionId=" + versionStr
				}
			}
			return url, nil
		}
	}

	// If any values are intrinsic functions, use Fn::Sub
	sub := map[string]interface{}{
		"Bucket": bucket,
		"Key":    key,
	}

	template := "https://${Bucket}.s3.amazonaws.com/${Key}"
	if version != nil {
		template += "?versionId=${Version}"
		sub["Version"] = version
	}

	return map[string]interface{}{
		"Fn::Sub": []interface{}{template, sub},
	}, nil
}
