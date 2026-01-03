// Package apigateway provides CloudFormation resource models for AWS API Gateway (REST API).
package apigateway

// Resource represents an AWS::ApiGateway::Resource CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-resource.html
type Resource struct {
	// ParentId is the ID of the parent resource.
	ParentId interface{} `json:"ParentId" yaml:"ParentId"`

	// PathPart is the final segment of this resource's path.
	PathPart interface{} `json:"PathPart" yaml:"PathPart"`

	// RestApiId is the ID of the RestApi resource in which to create this resource.
	RestApiId interface{} `json:"RestApiId" yaml:"RestApiId"`
}
