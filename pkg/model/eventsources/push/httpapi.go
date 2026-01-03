// Package push provides event source handlers for push-based Lambda integrations.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/apigatewayv2"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
	"github.com/lex00/aws-sam-translator-go/pkg/translator"
)

const (
	// ResourceTypeApi represents AWS::ApiGatewayV2::Api
	ResourceTypeApi = "AWS::ApiGatewayV2::Api"
	// ResourceTypeRoute represents AWS::ApiGatewayV2::Route
	ResourceTypeRoute = "AWS::ApiGatewayV2::Route"
	// ResourceTypeIntegration represents AWS::ApiGatewayV2::Integration
	ResourceTypeIntegration = "AWS::ApiGatewayV2::Integration"
	// ResourceTypeStage represents AWS::ApiGatewayV2::Stage
	ResourceTypeStage = "AWS::ApiGatewayV2::Stage"
)

// HttpApiEvent represents an HttpApi event source for AWS::Serverless::Function.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-httpapi.html
type HttpApiEvent struct {
	// ApiId is the identifier of an AWS::Serverless::HttpApi resource or AWS::ApiGatewayV2::Api.
	// If not specified, a default ServerlessHttpApi is created.
	ApiId interface{} `json:"ApiId,omitempty" yaml:"ApiId,omitempty"`

	// Auth specifies authorization configuration for this route.
	Auth *HttpApiAuth `json:"Auth,omitempty" yaml:"Auth,omitempty"`

	// Method is the HTTP method for which this function is invoked.
	// Valid values: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, ANY, or * (wildcard)
	// Default: ANY
	Method interface{} `json:"Method,omitempty" yaml:"Method,omitempty"`

	// Path is the URI path for which this function is invoked.
	// Must start with /.
	// Default: / (when ApiId is not specified)
	Path interface{} `json:"Path,omitempty" yaml:"Path,omitempty"`

	// PayloadFormatVersion specifies the format of the payload sent to the integration.
	// Valid values: 1.0, 2.0
	// Default: 2.0
	PayloadFormatVersion interface{} `json:"PayloadFormatVersion,omitempty" yaml:"PayloadFormatVersion,omitempty"`

	// RouteSettings specifies route-specific throttling and logging settings.
	RouteSettings *HttpApiRouteSettings `json:"RouteSettings,omitempty" yaml:"RouteSettings,omitempty"`

	// TimeoutInMillis is the custom timeout for the integration in milliseconds.
	// Valid range: 50-30000
	TimeoutInMillis interface{} `json:"TimeoutInMillis,omitempty" yaml:"TimeoutInMillis,omitempty"`
}

// HttpApiAuth represents authorization configuration for an HttpApi route.
type HttpApiAuth struct {
	// Authorizer is the name of the Lambda or JWT authorizer.
	// Special value NONE disables authorization for this route.
	// Can be AWS_IAM for IAM authorization.
	Authorizer interface{} `json:"Authorizer,omitempty" yaml:"Authorizer,omitempty"`

	// AuthorizationScopes is a list of authorization scopes for JWT authorizers.
	AuthorizationScopes []interface{} `json:"AuthorizationScopes,omitempty" yaml:"AuthorizationScopes,omitempty"`

	// InvokeRole is the ARN of the IAM role to use for integration credentials.
	// Special value CALLER_CREDENTIALS uses caller's credentials.
	InvokeRole interface{} `json:"InvokeRole,omitempty" yaml:"InvokeRole,omitempty"`
}

// HttpApiRouteSettings represents route-specific settings.
type HttpApiRouteSettings struct {
	// DataTraceEnabled specifies whether data trace logging is enabled.
	DataTraceEnabled interface{} `json:"DataTraceEnabled,omitempty" yaml:"DataTraceEnabled,omitempty"`

	// DetailedMetricsEnabled specifies whether detailed metrics are enabled.
	DetailedMetricsEnabled interface{} `json:"DetailedMetricsEnabled,omitempty" yaml:"DetailedMetricsEnabled,omitempty"`

	// LoggingLevel specifies the logging level.
	// Valid values: ERROR, INFO, OFF
	LoggingLevel interface{} `json:"LoggingLevel,omitempty" yaml:"LoggingLevel,omitempty"`

	// ThrottlingBurstLimit is the throttling burst limit.
	ThrottlingBurstLimit interface{} `json:"ThrottlingBurstLimit,omitempty" yaml:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit is the throttling rate limit.
	ThrottlingRateLimit interface{} `json:"ThrottlingRateLimit,omitempty" yaml:"ThrottlingRateLimit,omitempty"`
}

// HttpApiEventHandler handles the transformation of HttpApi event sources.
type HttpApiEventHandler struct {
	functionLogicalID string
	eventLogicalID    string
	event             *HttpApiEvent
}

// NewHttpApiEventHandler creates a new HttpApiEventHandler.
func NewHttpApiEventHandler(functionLogicalID, eventLogicalID string, event *HttpApiEvent) *HttpApiEventHandler {
	return &HttpApiEventHandler{
		functionLogicalID: functionLogicalID,
		eventLogicalID:    eventLogicalID,
		event:             event,
	}
}

// GenerateResources generates CloudFormation resources for the HttpApi event source.
// This includes:
// - AWS::ApiGatewayV2::Route (the route for this path/method)
// - AWS::ApiGatewayV2::Integration (Lambda integration)
// - AWS::Lambda::Permission (allows API Gateway to invoke the function)
// - Optionally, AWS::ApiGatewayV2::Api and AWS::ApiGatewayV2::Stage (if ApiId not specified)
func (h *HttpApiEventHandler) GenerateResources(functionArn interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Determine the API ID - either user-provided or implicit
	apiID := h.event.ApiId
	if apiID == nil {
		// If no ApiId specified, use the implicit ServerlessHttpApi
		apiID = map[string]interface{}{
			"Ref": "ServerlessHttpApi",
		}
	}

	// Determine path and method
	path := h.event.Path
	if path == nil {
		path = "/"
	}

	method := h.event.Method
	if method == nil {
		method = "ANY"
	}

	// Determine payload format version
	payloadFormatVersion := h.event.PayloadFormatVersion
	if payloadFormatVersion == nil {
		payloadFormatVersion = "2.0"
	}

	// Generate resource IDs
	idGenerator := translator.NewLogicalIDGenerator()
	integrationID := idGenerator.Generate(h.functionLogicalID, h.eventLogicalID, "Integration")
	routeID := idGenerator.Generate(h.functionLogicalID, h.eventLogicalID, "Route")
	permissionID := idGenerator.Generate(h.functionLogicalID, h.eventLogicalID, "Permission")

	// Create the Integration resource
	integration := &apigatewayv2.Integration{
		ApiId:                apiID,
		IntegrationType:      "AWS_PROXY",
		IntegrationUri:       functionArn,
		PayloadFormatVersion: payloadFormatVersion,
	}

	if h.event.TimeoutInMillis != nil {
		integration.TimeoutInMillis = h.event.TimeoutInMillis
	}

	// Handle credentials/invoke role
	if h.event.Auth != nil && h.event.Auth.InvokeRole != nil {
		integration.CredentialsArn = h.event.Auth.InvokeRole
	}

	// Create the Route resource
	routeKey := h.buildRouteKey(method, path)
	route := &apigatewayv2.Route{
		ApiId:    apiID,
		RouteKey: routeKey,
		Target: map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("integrations/${%s}", integrationID),
		},
	}

	// Handle authorization
	if h.event.Auth != nil {
		if err := h.applyAuthToRoute(route); err != nil {
			return nil, err
		}
	}

	// Create the Lambda Permission
	permission := lambda.NewAPIGatewayPermission(
		map[string]interface{}{"Ref": h.functionLogicalID},
		h.buildPermissionSourceArn(apiID),
	)

	// Convert to CloudFormation format
	resources[integrationID] = h.toCloudFormationResource(ResourceTypeIntegration, integration)
	resources[routeID] = h.toCloudFormationResource(ResourceTypeRoute, route)
	resources[permissionID] = permission.ToCloudFormation()

	return resources, nil
}

// buildRouteKey constructs the route key from method and path.
// Format: "METHOD /path" or "ANY /path" or "$default" for catch-all
func (h *HttpApiEventHandler) buildRouteKey(method, path interface{}) interface{} {
	// If both are simple strings, concatenate them
	methodStr, methodOK := method.(string)
	pathStr, pathOK := path.(string)

	if methodOK && pathOK {
		if pathStr == "/$default" {
			return "$default"
		}
		if methodStr == "*" || methodStr == "ANY" {
			methodStr = "ANY"
		}
		return fmt.Sprintf("%s %s", methodStr, pathStr)
	}

	// If dynamic values, use Fn::Sub
	return map[string]interface{}{
		"Fn::Sub": []interface{}{
			"${Method} ${Path}",
			map[string]interface{}{
				"Method": method,
				"Path":   path,
			},
		},
	}
}

// buildPermissionSourceArn constructs the source ARN for the Lambda permission.
func (h *HttpApiEventHandler) buildPermissionSourceArn(apiID interface{}) interface{} {
	// Format: arn:aws:execute-api:region:account-id:api-id/*
	return map[string]interface{}{
		"Fn::Sub": []interface{}{
			"arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${ApiId}/*",
			map[string]interface{}{
				"ApiId": apiID,
			},
		},
	}
}

// applyAuthToRoute applies authorization configuration to a route.
func (h *HttpApiEventHandler) applyAuthToRoute(route *apigatewayv2.Route) error {
	if h.event.Auth == nil {
		return nil
	}

	auth := h.event.Auth

	// Handle NONE - no authorization
	if auth.Authorizer != nil {
		authStr, ok := auth.Authorizer.(string)
		if ok && authStr == "NONE" {
			route.AuthorizationType = "NONE"
			return nil
		}

		// Handle AWS_IAM
		if ok && authStr == "AWS_IAM" {
			route.AuthorizationType = "AWS_IAM"
			return nil
		}

		// Handle JWT authorizer - the Authorizer is the authorizer name
		// The actual authorizer must be defined in the HttpApi resource
		route.AuthorizationType = "JWT"
		route.AuthorizerId = auth.Authorizer

		// Add authorization scopes if specified
		if len(auth.AuthorizationScopes) > 0 {
			route.AuthorizationScopes = auth.AuthorizationScopes
		}
	}

	return nil
}

// toCloudFormationResource wraps a resource with Type and Properties.
func (h *HttpApiEventHandler) toCloudFormationResource(resourceType string, resource interface{}) map[string]interface{} {
	cfnResource := make(map[string]interface{})
	cfnResource["Type"] = resourceType

	// Convert the resource to a map
	properties := make(map[string]interface{})

	switch r := resource.(type) {
	case *apigatewayv2.Integration:
		properties = h.integrationToMap(r)
	case *apigatewayv2.Route:
		properties = h.routeToMap(r)
	}

	cfnResource["Properties"] = properties
	return cfnResource
}

// integrationToMap converts Integration to a property map.
func (h *HttpApiEventHandler) integrationToMap(i *apigatewayv2.Integration) map[string]interface{} {
	m := make(map[string]interface{})

	m["ApiId"] = i.ApiId
	m["IntegrationType"] = i.IntegrationType

	if i.IntegrationUri != nil {
		m["IntegrationUri"] = i.IntegrationUri
	}
	if i.PayloadFormatVersion != nil {
		m["PayloadFormatVersion"] = i.PayloadFormatVersion
	}
	if i.TimeoutInMillis != nil {
		m["TimeoutInMillis"] = i.TimeoutInMillis
	}
	if i.CredentialsArn != nil {
		m["CredentialsArn"] = i.CredentialsArn
	}

	return m
}

// routeToMap converts Route to a property map.
func (h *HttpApiEventHandler) routeToMap(r *apigatewayv2.Route) map[string]interface{} {
	m := make(map[string]interface{})

	m["ApiId"] = r.ApiId
	m["RouteKey"] = r.RouteKey

	if r.Target != nil {
		m["Target"] = r.Target
	}
	if r.AuthorizationType != nil {
		m["AuthorizationType"] = r.AuthorizationType
	}
	if r.AuthorizerId != nil {
		m["AuthorizerId"] = r.AuthorizerId
	}
	if len(r.AuthorizationScopes) > 0 {
		m["AuthorizationScopes"] = r.AuthorizationScopes
	}

	return m
}
