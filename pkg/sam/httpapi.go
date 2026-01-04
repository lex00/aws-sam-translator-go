// Package sam provides SAM resource transformers.
package sam

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HttpApi represents an AWS::Serverless::HttpApi resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-httpapi.html
type HttpApi struct {
	// StageName is the name of the API stage. Default: $default
	StageName interface{} `json:"StageName,omitempty" yaml:"StageName,omitempty"`

	// Name is the name of the HTTP API resource.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Description is a description of the API.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DefinitionUri is the S3 URI or local path to an OpenAPI 3.0 definition.
	DefinitionUri interface{} `json:"DefinitionUri,omitempty" yaml:"DefinitionUri,omitempty"`

	// DefinitionBody is an inline OpenAPI 3.0 definition.
	DefinitionBody map[string]interface{} `json:"DefinitionBody,omitempty" yaml:"DefinitionBody,omitempty"`

	// CorsConfiguration specifies CORS configuration for the API.
	CorsConfiguration interface{} `json:"CorsConfiguration,omitempty" yaml:"CorsConfiguration,omitempty"`

	// Auth specifies authorization configuration.
	Auth *HttpApiAuth `json:"Auth,omitempty" yaml:"Auth,omitempty"`

	// AccessLogSettings specifies access logging configuration.
	AccessLogSettings *HttpApiAccessLogSettings `json:"AccessLogSettings,omitempty" yaml:"AccessLogSettings,omitempty"`

	// DefaultRouteSettings specifies default route settings.
	DefaultRouteSettings *HttpApiRouteSettings `json:"DefaultRouteSettings,omitempty" yaml:"DefaultRouteSettings,omitempty"`

	// RouteSettings is a map of route-specific settings.
	RouteSettings map[string]interface{} `json:"RouteSettings,omitempty" yaml:"RouteSettings,omitempty"`

	// StageVariables is a map of stage variables.
	StageVariables map[string]interface{} `json:"StageVariables,omitempty" yaml:"StageVariables,omitempty"`

	// Domain specifies custom domain configuration.
	Domain *HttpApiDomain `json:"Domain,omitempty" yaml:"Domain,omitempty"`

	// FailOnWarnings specifies whether to fail on import warnings.
	FailOnWarnings interface{} `json:"FailOnWarnings,omitempty" yaml:"FailOnWarnings,omitempty"`

	// DisableExecuteApiEndpoint disables the default execute-api endpoint.
	DisableExecuteApiEndpoint interface{} `json:"DisableExecuteApiEndpoint,omitempty" yaml:"DisableExecuteApiEndpoint,omitempty"`

	// Tags is a map of key-value pairs to apply to the API.
	Tags map[string]interface{} `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// HttpApiAuth specifies authorization configuration for HTTP API.
type HttpApiAuth struct {
	// DefaultAuthorizer is the default authorizer for all routes.
	DefaultAuthorizer string `json:"DefaultAuthorizer,omitempty" yaml:"DefaultAuthorizer,omitempty"`

	// Authorizers is a map of authorizer configurations.
	Authorizers map[string]interface{} `json:"Authorizers,omitempty" yaml:"Authorizers,omitempty"`

	// EnableIamAuthorizer enables IAM authorization.
	EnableIamAuthorizer bool `json:"EnableIamAuthorizer,omitempty" yaml:"EnableIamAuthorizer,omitempty"`
}

// HttpApiAccessLogSettings specifies access logging configuration.
type HttpApiAccessLogSettings struct {
	// DestinationArn is the ARN of the CloudWatch log group.
	DestinationArn interface{} `json:"DestinationArn,omitempty" yaml:"DestinationArn,omitempty"`

	// Format is the log format string.
	Format interface{} `json:"Format,omitempty" yaml:"Format,omitempty"`
}

// HttpApiRouteSettings specifies route settings for HTTP API.
type HttpApiRouteSettings struct {
	// ThrottlingBurstLimit is the throttling burst limit.
	ThrottlingBurstLimit interface{} `json:"ThrottlingBurstLimit,omitempty" yaml:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit is the throttling rate limit.
	ThrottlingRateLimit interface{} `json:"ThrottlingRateLimit,omitempty" yaml:"ThrottlingRateLimit,omitempty"`

	// DetailedMetricsEnabled enables detailed CloudWatch metrics.
	DetailedMetricsEnabled interface{} `json:"DetailedMetricsEnabled,omitempty" yaml:"DetailedMetricsEnabled,omitempty"`
}

// HttpApiDomain specifies custom domain configuration.
type HttpApiDomain struct {
	// DomainName is the custom domain name.
	DomainName interface{} `json:"DomainName,omitempty" yaml:"DomainName,omitempty"`

	// CertificateArn is the ARN of the certificate.
	CertificateArn interface{} `json:"CertificateArn,omitempty" yaml:"CertificateArn,omitempty"`

	// EndpointConfiguration is the endpoint type (REGIONAL).
	EndpointConfiguration interface{} `json:"EndpointConfiguration,omitempty" yaml:"EndpointConfiguration,omitempty"`

	// SecurityPolicy is the TLS version. Valid values: TLS_1_0, TLS_1_2
	SecurityPolicy interface{} `json:"SecurityPolicy,omitempty" yaml:"SecurityPolicy,omitempty"`

	// BasePath is a list of base path mappings.
	BasePath interface{} `json:"BasePath,omitempty" yaml:"BasePath,omitempty"`

	// MutualTlsAuthentication specifies mutual TLS configuration.
	MutualTlsAuthentication map[string]interface{} `json:"MutualTlsAuthentication,omitempty" yaml:"MutualTlsAuthentication,omitempty"`

	// Route53 specifies Route 53 configuration.
	Route53 map[string]interface{} `json:"Route53,omitempty" yaml:"Route53,omitempty"`
}

// HttpApiTransformer transforms AWS::Serverless::HttpApi to CloudFormation.
type HttpApiTransformer struct{}

// NewHttpApiTransformer creates a new HttpApiTransformer.
func NewHttpApiTransformer() *HttpApiTransformer {
	return &HttpApiTransformer{}
}

// Transform converts a SAM HttpApi to CloudFormation resources.
func (t *HttpApiTransformer) Transform(logicalID string, api *HttpApi, ctx *TransformContext) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Build the API Gateway V2 API resource
	apiProps, err := t.buildApiProperties(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build API properties: %w", err)
	}

	resources[logicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGatewayV2::Api",
		"Properties": apiProps,
	}

	// Build the Stage resource
	stageName := t.getStageName(api)
	stageLogicalID := logicalID + "Stage"
	stageProps := t.buildStageProperties(logicalID, api, stageName)

	resources[stageLogicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGatewayV2::Stage",
		"Properties": stageProps,
	}

	// Build authorizers if configured
	if api.Auth != nil && len(api.Auth.Authorizers) > 0 {
		authorizerResources, err := t.buildAuthorizers(logicalID, api.Auth)
		if err != nil {
			return nil, fmt.Errorf("failed to build authorizers: %w", err)
		}
		for k, v := range authorizerResources {
			resources[k] = v
		}
	}

	// Build custom domain resources if configured
	if api.Domain != nil {
		domainResources, err := t.buildDomainResources(logicalID, api.Domain, stageName)
		if err != nil {
			return nil, fmt.Errorf("failed to build domain resources: %w", err)
		}
		for k, v := range domainResources {
			resources[k] = v
		}
	}

	return resources, nil
}

// buildApiProperties builds the AWS::ApiGatewayV2::Api properties.
func (t *HttpApiTransformer) buildApiProperties(logicalID string, api *HttpApi) (map[string]interface{}, error) {
	props := make(map[string]interface{})

	// Protocol type is always HTTP for HttpApi
	props["ProtocolType"] = "HTTP"

	// Set Name
	if api.Name != nil {
		props["Name"] = api.Name
	} else {
		// Generate name from logical ID
		props["Name"] = logicalID
	}

	// Set Description
	if api.Description != nil {
		props["Description"] = api.Description
	}

	// Handle OpenAPI definition
	if api.DefinitionBody != nil {
		// Process and set Body
		body, err := t.processDefinitionBody(api.DefinitionBody)
		if err != nil {
			return nil, fmt.Errorf("failed to process definition body: %w", err)
		}
		props["Body"] = body
	} else if api.DefinitionUri != nil {
		// Set BodyS3Location
		s3Location, err := t.buildBodyS3Location(api.DefinitionUri)
		if err != nil {
			return nil, fmt.Errorf("failed to build S3 location: %w", err)
		}
		props["BodyS3Location"] = s3Location
	}

	// CORS configuration
	if api.CorsConfiguration != nil {
		corsConfig := t.buildCorsConfiguration(api.CorsConfiguration)
		if corsConfig != nil {
			props["CorsConfiguration"] = corsConfig
		}
	}

	// Fail on warnings
	if api.FailOnWarnings != nil {
		props["FailOnWarnings"] = api.FailOnWarnings
	}

	// Disable execute API endpoint
	if api.DisableExecuteApiEndpoint != nil {
		props["DisableExecuteApiEndpoint"] = api.DisableExecuteApiEndpoint
	}

	// Tags
	if len(api.Tags) > 0 {
		props["Tags"] = api.Tags
	}

	return props, nil
}

// getStageName returns the stage name, defaulting to "$default".
func (t *HttpApiTransformer) getStageName(api *HttpApi) interface{} {
	if api.StageName != nil {
		return api.StageName
	}
	return "$default"
}

// buildStageProperties builds the AWS::ApiGatewayV2::Stage properties.
func (t *HttpApiTransformer) buildStageProperties(apiLogicalID string, api *HttpApi, stageName interface{}) map[string]interface{} {
	props := map[string]interface{}{
		"ApiId":      map[string]interface{}{"Ref": apiLogicalID},
		"StageName":  stageName,
		"AutoDeploy": true,
	}

	// Access log settings
	if api.AccessLogSettings != nil {
		accessLogSettings := make(map[string]interface{})
		if api.AccessLogSettings.DestinationArn != nil {
			accessLogSettings["DestinationArn"] = api.AccessLogSettings.DestinationArn
		}
		if api.AccessLogSettings.Format != nil {
			accessLogSettings["Format"] = api.AccessLogSettings.Format
		} else {
			// Default format if destination is set but format is not
			if api.AccessLogSettings.DestinationArn != nil {
				accessLogSettings["Format"] = t.defaultAccessLogFormat()
			}
		}
		if len(accessLogSettings) > 0 {
			props["AccessLogSettings"] = accessLogSettings
		}
	}

	// Default route settings
	if api.DefaultRouteSettings != nil {
		defaultRouteSettings := make(map[string]interface{})
		if api.DefaultRouteSettings.ThrottlingBurstLimit != nil {
			defaultRouteSettings["ThrottlingBurstLimit"] = api.DefaultRouteSettings.ThrottlingBurstLimit
		}
		if api.DefaultRouteSettings.ThrottlingRateLimit != nil {
			defaultRouteSettings["ThrottlingRateLimit"] = api.DefaultRouteSettings.ThrottlingRateLimit
		}
		if api.DefaultRouteSettings.DetailedMetricsEnabled != nil {
			defaultRouteSettings["DetailedMetricsEnabled"] = api.DefaultRouteSettings.DetailedMetricsEnabled
		}
		if len(defaultRouteSettings) > 0 {
			props["DefaultRouteSettings"] = defaultRouteSettings
		}
	}

	// Route settings
	if len(api.RouteSettings) > 0 {
		props["RouteSettings"] = api.RouteSettings
	}

	// Stage variables
	if len(api.StageVariables) > 0 {
		props["StageVariables"] = api.StageVariables
	}

	// Tags
	if len(api.Tags) > 0 {
		props["Tags"] = api.Tags
	}

	return props
}

// defaultAccessLogFormat returns the default access log format for HTTP API.
func (t *HttpApiTransformer) defaultAccessLogFormat() string {
	return `{"requestId":"$context.requestId","ip":"$context.identity.sourceIp","requestTime":"$context.requestTime","httpMethod":"$context.httpMethod","routeKey":"$context.routeKey","status":"$context.status","protocol":"$context.protocol","responseLength":"$context.responseLength"}`
}

// processDefinitionBody processes the OpenAPI definition body.
func (t *HttpApiTransformer) processDefinitionBody(body map[string]interface{}) (interface{}, error) {
	// Check if it contains intrinsic functions that need to be preserved
	if t.containsIntrinsics(body) {
		return body, nil
	}

	// Convert to JSON string with Fn::Join for proper CloudFormation handling
	jsonBytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal definition body: %w", err)
	}

	return string(jsonBytes), nil
}

// containsIntrinsics checks if the map contains CloudFormation intrinsic functions.
func (t *HttpApiTransformer) containsIntrinsics(m map[string]interface{}) bool {
	for k, v := range m {
		// Check for common intrinsic functions
		if strings.HasPrefix(k, "Fn::") || k == "Ref" || k == "Condition" {
			return true
		}
		// Recursively check nested maps
		if nested, ok := v.(map[string]interface{}); ok {
			if t.containsIntrinsics(nested) {
				return true
			}
		}
		// Check arrays
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					if t.containsIntrinsics(nestedMap) {
						return true
					}
				}
			}
		}
	}
	return false
}

// buildBodyS3Location builds the BodyS3Location from DefinitionUri.
func (t *HttpApiTransformer) buildBodyS3Location(definitionUri interface{}) (map[string]interface{}, error) {
	switch uri := definitionUri.(type) {
	case string:
		return t.parseS3Uri(uri)
	case map[string]interface{}:
		s3Location := make(map[string]interface{})
		if bucket, ok := uri["Bucket"]; ok {
			s3Location["Bucket"] = bucket
		}
		if key, ok := uri["Key"]; ok {
			s3Location["Key"] = key
		}
		if version, ok := uri["Version"]; ok {
			s3Location["Version"] = version
		}
		return s3Location, nil
	default:
		return nil, fmt.Errorf("unsupported DefinitionUri type: %T", definitionUri)
	}
}

// parseS3Uri parses an S3 URI string into an S3 location map.
func (t *HttpApiTransformer) parseS3Uri(uri string) (map[string]interface{}, error) {
	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("DefinitionUri must be an S3 URI (s3://...): %s", uri)
	}

	path := strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid S3 URI (missing key): %s", uri)
	}

	return map[string]interface{}{
		"Bucket": parts[0],
		"Key":    parts[1],
	}, nil
}

// buildCorsConfiguration builds the CORS configuration.
func (t *HttpApiTransformer) buildCorsConfiguration(corsConfig interface{}) map[string]interface{} {
	switch cors := corsConfig.(type) {
	case bool:
		if cors {
			// Simple CORS - allow all origins
			return map[string]interface{}{
				"AllowOrigins": []interface{}{"*"},
				"AllowMethods": []interface{}{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
				"AllowHeaders": []interface{}{"Content-Type", "X-Amz-Date", "Authorization", "X-Api-Key", "X-Amz-Security-Token"},
			}
		}
		return nil
	case string:
		// Single origin
		return map[string]interface{}{
			"AllowOrigins": []interface{}{cors},
		}
	case map[string]interface{}:
		// Full CORS configuration
		result := make(map[string]interface{})
		if origins, ok := cors["AllowOrigins"]; ok {
			result["AllowOrigins"] = origins
		}
		if methods, ok := cors["AllowMethods"]; ok {
			result["AllowMethods"] = methods
		}
		if headers, ok := cors["AllowHeaders"]; ok {
			result["AllowHeaders"] = headers
		}
		if exposeHeaders, ok := cors["ExposeHeaders"]; ok {
			result["ExposeHeaders"] = exposeHeaders
		}
		if credentials, ok := cors["AllowCredentials"]; ok {
			result["AllowCredentials"] = credentials
		}
		if maxAge, ok := cors["MaxAge"]; ok {
			result["MaxAge"] = maxAge
		}
		return result
	default:
		return nil
	}
}

// buildAuthorizers builds authorizer resources.
func (t *HttpApiTransformer) buildAuthorizers(apiLogicalID string, auth *HttpApiAuth) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	for authName, authConfig := range auth.Authorizers {
		authLogicalID := apiLogicalID + authName + "Authorizer"

		authMap, ok := authConfig.(map[string]interface{})
		if !ok {
			continue
		}

		authProps := map[string]interface{}{
			"ApiId": map[string]interface{}{"Ref": apiLogicalID},
			"Name":  authName,
		}

		// Determine authorizer type
		if jwtConfig, hasJwt := authMap["JwtConfiguration"]; hasJwt {
			authProps["AuthorizerType"] = "JWT"
			authProps["JwtConfiguration"] = jwtConfig

			// Identity source (defaults to Authorization header)
			if identitySource, ok := authMap["IdentitySource"]; ok {
				authProps["IdentitySource"] = identitySource
			} else {
				authProps["IdentitySource"] = []interface{}{"$request.header.Authorization"}
			}
		} else if _, hasLambda := authMap["FunctionArn"]; hasLambda {
			authProps["AuthorizerType"] = "REQUEST"

			// Build Lambda authorizer URI
			if functionArn, ok := authMap["FunctionArn"]; ok {
				authProps["AuthorizerUri"] = t.buildLambdaAuthorizerUri(functionArn)
			}

			// Authorizer payload format version
			if payloadVersion, ok := authMap["AuthorizerPayloadFormatVersion"]; ok {
				authProps["AuthorizerPayloadFormatVersion"] = payloadVersion
			} else {
				authProps["AuthorizerPayloadFormatVersion"] = "2.0"
			}

			// Enable simple responses
			if enableSimple, ok := authMap["EnableSimpleResponses"]; ok {
				authProps["EnableSimpleResponses"] = enableSimple
			}

			// Identity source
			if identitySource, ok := authMap["IdentitySource"]; ok {
				authProps["IdentitySource"] = identitySource
			}

			// Result TTL
			if ttl, ok := authMap["AuthorizerResultTtlInSeconds"]; ok {
				authProps["AuthorizerResultTtlInSeconds"] = ttl
			}
		}

		resources[authLogicalID] = map[string]interface{}{
			"Type":       "AWS::ApiGatewayV2::Authorizer",
			"Properties": authProps,
		}
	}

	return resources, nil
}

// buildLambdaAuthorizerUri builds the authorizer URI for a Lambda function.
func (t *HttpApiTransformer) buildLambdaAuthorizerUri(functionArn interface{}) interface{} {
	switch arn := functionArn.(type) {
	case string:
		return map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/%s/invocations", arn),
		}
	case map[string]interface{}:
		if _, hasRef := arn["Ref"]; hasRef {
			return map[string]interface{}{
				"Fn::Sub": []interface{}{
					"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
					map[string]interface{}{
						"FunctionArn": map[string]interface{}{
							"Fn::GetAtt": []interface{}{arn["Ref"], "Arn"},
						},
					},
				},
			}
		}
		if getAtt, hasGetAtt := arn["Fn::GetAtt"]; hasGetAtt {
			return map[string]interface{}{
				"Fn::Sub": []interface{}{
					"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
					map[string]interface{}{
						"FunctionArn": map[string]interface{}{"Fn::GetAtt": getAtt},
					},
				},
			}
		}
	}
	return functionArn
}

// buildDomainResources builds custom domain resources.
func (t *HttpApiTransformer) buildDomainResources(apiLogicalID string, domain *HttpApiDomain, stageName interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	if domain.DomainName == nil {
		return resources, nil
	}

	// Create Domain Name resource
	domainLogicalID := apiLogicalID + "DomainName"
	domainProps := map[string]interface{}{
		"DomainName": domain.DomainName,
	}

	// Domain name configurations
	domainNameConfig := map[string]interface{}{}

	if domain.CertificateArn != nil {
		domainNameConfig["CertificateArn"] = domain.CertificateArn
	}

	if domain.EndpointConfiguration != nil {
		domainNameConfig["EndpointType"] = domain.EndpointConfiguration
	} else {
		domainNameConfig["EndpointType"] = "REGIONAL"
	}

	if domain.SecurityPolicy != nil {
		domainNameConfig["SecurityPolicy"] = domain.SecurityPolicy
	}

	domainProps["DomainNameConfigurations"] = []interface{}{domainNameConfig}

	// Mutual TLS
	if domain.MutualTlsAuthentication != nil {
		domainProps["MutualTlsAuthentication"] = domain.MutualTlsAuthentication
	}

	resources[domainLogicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGatewayV2::DomainName",
		"Properties": domainProps,
	}

	// Create API Mapping
	mappingLogicalID := apiLogicalID + "ApiMapping"
	mappingProps := map[string]interface{}{
		"ApiId":      map[string]interface{}{"Ref": apiLogicalID},
		"DomainName": map[string]interface{}{"Ref": domainLogicalID},
		"Stage":      stageName,
	}

	// Base path
	if domain.BasePath != nil {
		switch bp := domain.BasePath.(type) {
		case string:
			if bp != "" && bp != "/" {
				mappingProps["ApiMappingKey"] = bp
			}
		case []interface{}:
			// Multiple base paths - create multiple mappings
			for i, path := range bp {
				if pathStr, ok := path.(string); ok && pathStr != "" && pathStr != "/" {
					mappingID := fmt.Sprintf("%s%d", mappingLogicalID, i)
					resources[mappingID] = map[string]interface{}{
						"Type": "AWS::ApiGatewayV2::ApiMapping",
						"Properties": map[string]interface{}{
							"ApiId":         map[string]interface{}{"Ref": apiLogicalID},
							"DomainName":    map[string]interface{}{"Ref": domainLogicalID},
							"Stage":         stageName,
							"ApiMappingKey": pathStr,
						},
						"DependsOn": []interface{}{domainLogicalID, apiLogicalID + "Stage"},
					}
				}
			}
			// Return early if we created multiple mappings
			if len(bp) > 0 {
				return resources, nil
			}
		}
	}

	resources[mappingLogicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGatewayV2::ApiMapping",
		"Properties": mappingProps,
		"DependsOn":  []interface{}{domainLogicalID, apiLogicalID + "Stage"},
	}

	// Route 53 record if configured
	if domain.Route53 != nil {
		route53Resources, err := t.buildRoute53Resources(apiLogicalID, domain)
		if err != nil {
			return nil, err
		}
		for k, v := range route53Resources {
			resources[k] = v
		}
	}

	return resources, nil
}

// buildRoute53Resources builds Route 53 record set resources.
func (t *HttpApiTransformer) buildRoute53Resources(apiLogicalID string, domain *HttpApiDomain) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	route53Config := domain.Route53
	domainLogicalID := apiLogicalID + "DomainName"
	recordLogicalID := apiLogicalID + "RecordSet"

	recordProps := map[string]interface{}{
		"Name": domain.DomainName,
		"Type": "A",
		"AliasTarget": map[string]interface{}{
			"DNSName": map[string]interface{}{
				"Fn::GetAtt": []interface{}{domainLogicalID, "RegionalDomainName"},
			},
			"HostedZoneId": map[string]interface{}{
				"Fn::GetAtt": []interface{}{domainLogicalID, "RegionalHostedZoneId"},
			},
		},
	}

	// Hosted zone ID or name
	if hostedZoneId, ok := route53Config["HostedZoneId"]; ok {
		recordProps["HostedZoneId"] = hostedZoneId
	} else if hostedZoneName, ok := route53Config["HostedZoneName"]; ok {
		recordProps["HostedZoneName"] = hostedZoneName
	}

	// Set TTL for alias record (not applicable but some may specify)
	if evaluateTargetHealth, ok := route53Config["EvaluateTargetHealth"]; ok {
		if aliasTarget, ok := recordProps["AliasTarget"].(map[string]interface{}); ok {
			aliasTarget["EvaluateTargetHealth"] = evaluateTargetHealth
		}
	}

	resources[recordLogicalID] = map[string]interface{}{
		"Type":       "AWS::Route53::RecordSet",
		"Properties": recordProps,
	}

	// Create AAAA record for IPv6 if specified
	if ipv6, ok := route53Config["IpV6"].(bool); ok && ipv6 {
		ipv6RecordProps := make(map[string]interface{})
		for k, v := range recordProps {
			ipv6RecordProps[k] = v
		}
		ipv6RecordProps["Type"] = "AAAA"

		resources[recordLogicalID+"V6"] = map[string]interface{}{
			"Type":       "AWS::Route53::RecordSet",
			"Properties": ipv6RecordProps,
		}
	}

	return resources, nil
}
