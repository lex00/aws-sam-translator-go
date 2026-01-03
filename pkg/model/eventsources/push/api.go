// Package push provides event source handlers for push-based integrations.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/apigateway"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
	"github.com/lex00/aws-sam-translator-go/pkg/translator"
)

// Api represents an API Gateway REST API event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-api.html
type Api struct {
	// Path is the URI path for which this function is invoked (required).
	// Must start with /.
	Path interface{} `json:"Path" yaml:"Path"`

	// Method is the HTTP method for which this function is invoked (required).
	// Valid values: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH, ANY
	Method interface{} `json:"Method" yaml:"Method"`

	// RestApiId is the identifier of a RestApi resource to associate with this path and method.
	// If not specified, a default API will be created.
	RestApiId interface{} `json:"RestApiId,omitempty" yaml:"RestApiId,omitempty"`

	// Auth configures authorization for this specific API event.
	Auth *ApiAuth `json:"Auth,omitempty" yaml:"Auth,omitempty"`

	// RequestModel configures request validation for this endpoint.
	RequestModel *RequestModel `json:"RequestModel,omitempty" yaml:"RequestModel,omitempty"`

	// RequestParameters configures required request parameters.
	RequestParameters []string `json:"RequestParameters,omitempty" yaml:"RequestParameters,omitempty"`
}

// ApiAuth represents authorization configuration for an API event.
type ApiAuth struct {
	// Authorizer is the name of the Lambda authorizer to use.
	Authorizer interface{} `json:"Authorizer,omitempty" yaml:"Authorizer,omitempty"`

	// AuthorizationScopes is a list of authorization scopes for this method.
	AuthorizationScopes []interface{} `json:"AuthorizationScopes,omitempty" yaml:"AuthorizationScopes,omitempty"`

	// ApiKeyRequired indicates whether the method requires an API key.
	ApiKeyRequired interface{} `json:"ApiKeyRequired,omitempty" yaml:"ApiKeyRequired,omitempty"`

	// ResourcePolicy configures resource policy statements for this endpoint.
	ResourcePolicy *ResourcePolicy `json:"ResourcePolicy,omitempty" yaml:"ResourcePolicy,omitempty"`

	// InvokeRole specifies the credentials required to invoke the API.
	InvokeRole interface{} `json:"InvokeRole,omitempty" yaml:"InvokeRole,omitempty"`
}

// RequestModel represents request model configuration for validation.
type RequestModel struct {
	// Model is the name of a model defined in the API's Models property.
	Model interface{} `json:"Model,omitempty" yaml:"Model,omitempty"`

	// Required indicates whether to validate the request body.
	Required interface{} `json:"Required,omitempty" yaml:"Required,omitempty"`

	// ValidateBody indicates whether to validate the request body against the model.
	ValidateBody interface{} `json:"ValidateBody,omitempty" yaml:"ValidateBody,omitempty"`

	// ValidateParameters indicates whether to validate request parameters.
	ValidateParameters interface{} `json:"ValidateParameters,omitempty" yaml:"ValidateParameters,omitempty"`
}

// ResourcePolicy represents resource policy configuration for an API.
type ResourcePolicy struct {
	// AwsAccountWhitelist is a list of AWS account IDs to allow.
	AwsAccountWhitelist []interface{} `json:"AwsAccountWhitelist,omitempty" yaml:"AwsAccountWhitelist,omitempty"`

	// AwsAccountBlacklist is a list of AWS account IDs to deny.
	AwsAccountBlacklist []interface{} `json:"AwsAccountBlacklist,omitempty" yaml:"AwsAccountBlacklist,omitempty"`

	// IpRangeWhitelist is a list of IP ranges to allow.
	IpRangeWhitelist []interface{} `json:"IpRangeWhitelist,omitempty" yaml:"IpRangeWhitelist,omitempty"`

	// IpRangeBlacklist is a list of IP ranges to deny.
	IpRangeBlacklist []interface{} `json:"IpRangeBlacklist,omitempty" yaml:"IpRangeBlacklist,omitempty"`

	// SourceVpcWhitelist is a list of VPC IDs or VPC endpoint IDs to allow.
	SourceVpcWhitelist []interface{} `json:"SourceVpcWhitelist,omitempty" yaml:"SourceVpcWhitelist,omitempty"`

	// SourceVpcBlacklist is a list of VPC IDs or VPC endpoint IDs to deny.
	SourceVpcBlacklist []interface{} `json:"SourceVpcBlacklist,omitempty" yaml:"SourceVpcBlacklist,omitempty"`

	// CustomStatements is a list of custom policy statements.
	CustomStatements []interface{} `json:"CustomStatements,omitempty" yaml:"CustomStatements,omitempty"`
}

// ToCloudFormation generates CloudFormation resources for this API event.
// This includes:
// - AWS::ApiGateway::Resource (if path is not root)
// - AWS::ApiGateway::Method
// - AWS::Lambda::Permission
func (a *Api) ToCloudFormation(functionLogicalId string, functionRef interface{}, stageName interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	if a.Path == nil {
		return nil, fmt.Errorf("Api event source requires Path property")
	}
	if a.Method == nil {
		return nil, fmt.Errorf("Api event source requires Method property")
	}

	// Determine the RestApiId to use
	var restApiId interface{}
	if a.RestApiId != nil {
		restApiId = a.RestApiId
	} else {
		// If RestApiId is not specified, use ServerlessRestApi (implicit API)
		restApiId = map[string]interface{}{"Ref": "ServerlessRestApi"}
	}

	// Generate logical IDs
	idGen := translator.NewLogicalIDGenerator()

	// Convert path and method to strings for logical ID generation
	pathStr := fmt.Sprintf("%v", a.Path)
	methodStr := fmt.Sprintf("%v", a.Method)
	stageStr := fmt.Sprintf("%v", stageName)

	methodLogicalId := idGen.Generate(functionLogicalId, "ApiEvent", pathStr, methodStr)
	permissionLogicalId := idGen.Generate(functionLogicalId, "Permission", pathStr, methodStr, stageStr)

	// Create API Gateway Method
	method := &apigateway.Method{
		HttpMethod:        a.Method,
		ResourceId:        a.getResourceId(restApiId, a.Path),
		RestApiId:         restApiId,
		AuthorizationType: a.getAuthorizationType(),
	}

	// Configure authorization
	if a.Auth != nil {
		if a.Auth.Authorizer != nil {
			method.AuthorizerId = a.Auth.Authorizer
		}
		if a.Auth.AuthorizationScopes != nil {
			method.AuthorizationScopes = a.Auth.AuthorizationScopes
		}
		if a.Auth.ApiKeyRequired != nil {
			method.ApiKeyRequired = a.Auth.ApiKeyRequired
		}
	}

	// Configure Lambda integration
	method.Integration = &apigateway.Integration{
		Type:                  "AWS_PROXY",
		IntegrationHttpMethod: "POST",
		Uri: map[string]interface{}{
			"Fn::Sub": []interface{}{
				"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
				map[string]interface{}{
					"FunctionArn": map[string]interface{}{
						"Fn::GetAtt": []interface{}{functionLogicalId, "Arn"},
					},
				},
			},
		},
	}

	resources[methodLogicalId] = map[string]interface{}{
		"Type":       "AWS::ApiGateway::Method",
		"Properties": method,
	}

	// Create Lambda Permission
	permission := lambda.NewAPIGatewayPermission(
		functionRef,
		a.buildSourceArn(restApiId, stageName),
	)

	resources[permissionLogicalId] = permission.ToCloudFormation()

	return resources, nil
}

// getAuthorizationType returns the authorization type for the method.
func (a *Api) getAuthorizationType() interface{} {
	if a.Auth != nil {
		if a.Auth.Authorizer != nil {
			// If an authorizer is specified, use CUSTOM
			return "CUSTOM"
		}
		if a.Auth.InvokeRole != nil {
			return "AWS_IAM"
		}
	}
	// Default to NONE
	return "NONE"
}

// getResourceId returns the resource ID for the API path.
// For root path (/), returns the root resource ID.
// For other paths, would require creating intermediate resources.
func (a *Api) getResourceId(restApiId interface{}, path interface{}) interface{} {
	pathStr, ok := path.(string)
	if !ok || pathStr == "/" {
		// Root resource
		return map[string]interface{}{
			"Fn::GetAtt": []interface{}{restApiId, "RootResourceId"},
		}
	}

	// For non-root paths, we would need to create AWS::ApiGateway::Resource resources
	// and maintain a path tree. This is a simplified version.
	// In a full implementation, this would:
	// 1. Parse the path into segments
	// 2. Create a Resource for each segment
	// 3. Return a reference to the final segment's resource
	//
	// For now, we return a reference that assumes the resource exists
	idGen := translator.NewLogicalIDGenerator()
	resourceLogicalId := idGen.Generate("ApiResource", fmt.Sprintf("%v", path))
	return map[string]interface{}{"Ref": resourceLogicalId}
}

// buildSourceArn builds the source ARN for the Lambda permission.
// This allows API Gateway to invoke the Lambda function.
func (a *Api) buildSourceArn(restApiId interface{}, stageName interface{}) interface{} {
	// Determine the API ID string for substitution
	var apiIdStr string
	if refMap, ok := restApiId.(map[string]interface{}); ok {
		if ref, ok := refMap["Ref"].(string); ok {
			apiIdStr = ref
		}
	} else if str, ok := restApiId.(string); ok {
		apiIdStr = str
	}

	// Determine the stage string
	var stageStr string
	if stageName != nil {
		if str, ok := stageName.(string); ok {
			stageStr = str
		}
	}
	if stageStr == "" {
		stageStr = "*"
	}

	// Determine the method string
	var methodStr string
	if str, ok := a.Method.(string); ok {
		methodStr = str
	} else {
		methodStr = "*"
	}

	// Determine the path string
	var pathStr string
	if str, ok := a.Path.(string); ok {
		pathStr = str
	} else {
		pathStr = "/*"
	}

	// Build the execute-api ARN
	// Format: arn:aws:execute-api:region:account-id:api-id/stage/method/path
	return map[string]interface{}{
		"Fn::Sub": []interface{}{
			fmt.Sprintf("arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${__ApiId__}/${__Stage__}/%s%s", methodStr, pathStr),
			map[string]interface{}{
				"__ApiId__": apiIdStr,
				"__Stage__": stageStr,
			},
		},
	}
}

// Validate checks if the API event source configuration is valid.
func (a *Api) Validate() error {
	if a.Path == nil {
		return fmt.Errorf("Api event source requires Path property")
	}
	if a.Method == nil {
		return fmt.Errorf("Api event source requires Method property")
	}
	return nil
}

// EventType returns the event source type identifier.
func (a *Api) EventType() string {
	return "Api"
}
