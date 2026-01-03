// Package apigateway provides CloudFormation resource models for AWS API Gateway (REST API).
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
package apigateway

// RestApi represents an AWS::ApiGateway::RestApi CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
type RestApi struct {
	// ApiKeySourceType specifies the source of the API key for metering requests.
	// Valid values: HEADER, AUTHORIZER
	ApiKeySourceType interface{} `json:"ApiKeySourceType,omitempty" yaml:"ApiKeySourceType,omitempty"`

	// BinaryMediaTypes is the list of binary media types that are supported by the RestApi.
	BinaryMediaTypes []interface{} `json:"BinaryMediaTypes,omitempty" yaml:"BinaryMediaTypes,omitempty"`

	// Body is an OpenAPI specification that defines a set of RESTful APIs in JSON or YAML format.
	Body interface{} `json:"Body,omitempty" yaml:"Body,omitempty"`

	// BodyS3Location specifies the Amazon S3 location of the OpenAPI file that defines the API.
	BodyS3Location *S3Location `json:"BodyS3Location,omitempty" yaml:"BodyS3Location,omitempty"`

	// CloneFrom specifies the ID of the RestApi that you want to clone.
	CloneFrom interface{} `json:"CloneFrom,omitempty" yaml:"CloneFrom,omitempty"`

	// Description is a description of the RestApi resource.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DisableExecuteApiEndpoint specifies whether clients can invoke your API using the default execute-api endpoint.
	DisableExecuteApiEndpoint interface{} `json:"DisableExecuteApiEndpoint,omitempty" yaml:"DisableExecuteApiEndpoint,omitempty"`

	// EndpointConfiguration specifies the endpoint configuration for the REST API.
	EndpointConfiguration *EndpointConfiguration `json:"EndpointConfiguration,omitempty" yaml:"EndpointConfiguration,omitempty"`

	// FailOnWarnings indicates whether to rollback the resource if a warning occurs during API creation.
	FailOnWarnings interface{} `json:"FailOnWarnings,omitempty" yaml:"FailOnWarnings,omitempty"`

	// MinimumCompressionSize specifies the minimum compression size (in bytes) for responses.
	MinimumCompressionSize interface{} `json:"MinimumCompressionSize,omitempty" yaml:"MinimumCompressionSize,omitempty"`

	// Mode specifies how API Gateway handles resource updates.
	// Valid values: overwrite, merge
	Mode interface{} `json:"Mode,omitempty" yaml:"Mode,omitempty"`

	// Name is the name of the RestApi resource.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Parameters specifies custom header parameters for the request.
	Parameters map[string]interface{} `json:"Parameters,omitempty" yaml:"Parameters,omitempty"`

	// Policy is a policy document that contains permissions to add to the specified API.
	Policy interface{} `json:"Policy,omitempty" yaml:"Policy,omitempty"`

	// Tags is a map of key-value pairs to associate with the RestApi.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// S3Location represents an Amazon S3 location for an OpenAPI specification.
type S3Location struct {
	// Bucket is the name of the S3 bucket where the OpenAPI file is stored.
	Bucket interface{} `json:"Bucket,omitempty" yaml:"Bucket,omitempty"`

	// ETag is the ETag of the S3 object.
	ETag interface{} `json:"ETag,omitempty" yaml:"ETag,omitempty"`

	// Key is the name of the S3 object that represents the OpenAPI file.
	Key interface{} `json:"Key,omitempty" yaml:"Key,omitempty"`

	// Version is the version ID of the S3 object.
	Version interface{} `json:"Version,omitempty" yaml:"Version,omitempty"`
}

// EndpointConfiguration represents the endpoint configuration for a REST API.
type EndpointConfiguration struct {
	// Types is a list of endpoint types for the REST API.
	// Valid values: EDGE, REGIONAL, PRIVATE
	Types []interface{} `json:"Types,omitempty" yaml:"Types,omitempty"`

	// VpcEndpointIds is a list of VPC endpoint IDs for private APIs.
	VpcEndpointIds []interface{} `json:"VpcEndpointIds,omitempty" yaml:"VpcEndpointIds,omitempty"`
}

// Tag represents a key-value pair for tagging resources.
type Tag struct {
	// Key is the key of the tag.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the value of the tag.
	Value interface{} `json:"Value" yaml:"Value"`
}
