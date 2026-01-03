// Package iam provides IAM CloudFormation resource models for SAM transformation.
package iam

import (
	"fmt"
)

// Role represents an AWS::IAM::Role CloudFormation resource.
type Role struct {
	// RoleName is the name of the role. Optional - CloudFormation generates if not specified.
	RoleName interface{} `json:"RoleName,omitempty"`

	// AssumeRolePolicyDocument is the trust policy that is associated with this role.
	AssumeRolePolicyDocument *PolicyDocument `json:"AssumeRolePolicyDocument"`

	// Description provides a description of the role.
	Description string `json:"Description,omitempty"`

	// Path is the path to the role.
	Path string `json:"Path,omitempty"`

	// MaxSessionDuration is the maximum session duration (in seconds).
	MaxSessionDuration int `json:"MaxSessionDuration,omitempty"`

	// PermissionsBoundary is the ARN of the policy used to set permissions boundary.
	PermissionsBoundary interface{} `json:"PermissionsBoundary,omitempty"`

	// ManagedPolicyArns is a list of ARNs of managed policies to attach.
	ManagedPolicyArns []interface{} `json:"ManagedPolicyArns,omitempty"`

	// Policies is a list of inline policies embedded in the role.
	Policies []InlinePolicy `json:"Policies,omitempty"`

	// Tags is a list of tags for the role.
	Tags []Tag `json:"Tags,omitempty"`
}

// InlinePolicy represents an inline policy embedded in a role.
type InlinePolicy struct {
	// PolicyName is the name of the policy.
	PolicyName string `json:"PolicyName"`

	// PolicyDocument is the policy document.
	PolicyDocument *PolicyDocument `json:"PolicyDocument"`
}

// Tag represents a resource tag.
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// NewRole creates a new Role with the specified assume role policy.
func NewRole(assumeRolePolicy *PolicyDocument) *Role {
	return &Role{
		AssumeRolePolicyDocument: assumeRolePolicy,
	}
}

// NewRoleWithName creates a new Role with a name and assume role policy.
func NewRoleWithName(name interface{}, assumeRolePolicy *PolicyDocument) *Role {
	return &Role{
		RoleName:                 name,
		AssumeRolePolicyDocument: assumeRolePolicy,
	}
}

// WithPath sets the path for the role.
func (r *Role) WithPath(path string) *Role {
	r.Path = path
	return r
}

// WithDescription sets the description for the role.
func (r *Role) WithDescription(description string) *Role {
	r.Description = description
	return r
}

// WithMaxSessionDuration sets the max session duration for the role.
func (r *Role) WithMaxSessionDuration(duration int) *Role {
	r.MaxSessionDuration = duration
	return r
}

// WithPermissionsBoundary sets the permissions boundary for the role.
func (r *Role) WithPermissionsBoundary(arn interface{}) *Role {
	r.PermissionsBoundary = arn
	return r
}

// AddManagedPolicyArn adds a managed policy ARN to the role.
func (r *Role) AddManagedPolicyArn(arn interface{}) *Role {
	r.ManagedPolicyArns = append(r.ManagedPolicyArns, arn)
	return r
}

// AddManagedPolicyArns adds multiple managed policy ARNs to the role.
func (r *Role) AddManagedPolicyArns(arns []interface{}) *Role {
	r.ManagedPolicyArns = append(r.ManagedPolicyArns, arns...)
	return r
}

// AddInlinePolicy adds an inline policy to the role.
func (r *Role) AddInlinePolicy(name string, document *PolicyDocument) *Role {
	r.Policies = append(r.Policies, InlinePolicy{
		PolicyName:     name,
		PolicyDocument: document,
	})
	return r
}

// AddTag adds a tag to the role.
func (r *Role) AddTag(key, value string) *Role {
	r.Tags = append(r.Tags, Tag{Key: key, Value: value})
	return r
}

// ToCloudFormation converts the role to CloudFormation resource properties.
func (r *Role) ToCloudFormation() map[string]interface{} {
	props := make(map[string]interface{})

	if r.RoleName != nil {
		props["RoleName"] = r.RoleName
	}

	if r.AssumeRolePolicyDocument != nil {
		props["AssumeRolePolicyDocument"] = r.AssumeRolePolicyDocument.ToMap()
	}

	if r.Description != "" {
		props["Description"] = r.Description
	}

	if r.Path != "" {
		props["Path"] = r.Path
	}

	if r.MaxSessionDuration > 0 {
		props["MaxSessionDuration"] = r.MaxSessionDuration
	}

	if r.PermissionsBoundary != nil {
		props["PermissionsBoundary"] = r.PermissionsBoundary
	}

	if len(r.ManagedPolicyArns) > 0 {
		props["ManagedPolicyArns"] = r.ManagedPolicyArns
	}

	if len(r.Policies) > 0 {
		policies := make([]map[string]interface{}, len(r.Policies))
		for i, p := range r.Policies {
			policies[i] = map[string]interface{}{
				"PolicyName":     p.PolicyName,
				"PolicyDocument": p.PolicyDocument.ToMap(),
			}
		}
		props["Policies"] = policies
	}

	if len(r.Tags) > 0 {
		tags := make([]map[string]interface{}, len(r.Tags))
		for i, t := range r.Tags {
			tags[i] = map[string]interface{}{
				"Key":   t.Key,
				"Value": t.Value,
			}
		}
		props["Tags"] = tags
	}

	return props
}

// ToResource converts the role to a complete CloudFormation resource.
func (r *Role) ToResource() map[string]interface{} {
	return map[string]interface{}{
		"Type":       "AWS::IAM::Role",
		"Properties": r.ToCloudFormation(),
	}
}

// Validate validates the role configuration.
func (r *Role) Validate() error {
	if r.AssumeRolePolicyDocument == nil {
		return fmt.Errorf("AssumeRolePolicyDocument is required for IAM::Role")
	}

	if err := r.AssumeRolePolicyDocument.Validate(); err != nil {
		return fmt.Errorf("invalid AssumeRolePolicyDocument: %w", err)
	}

	for i, p := range r.Policies {
		if p.PolicyName == "" {
			return fmt.Errorf("inline policy at index %d is missing PolicyName", i)
		}
		if p.PolicyDocument == nil {
			return fmt.Errorf("inline policy '%s' is missing PolicyDocument", p.PolicyName)
		}
		if err := p.PolicyDocument.Validate(); err != nil {
			return fmt.Errorf("invalid PolicyDocument for inline policy '%s': %w", p.PolicyName, err)
		}
	}

	if r.MaxSessionDuration != 0 && (r.MaxSessionDuration < 3600 || r.MaxSessionDuration > 43200) {
		return fmt.Errorf("MaxSessionDuration must be between 3600 and 43200 seconds")
	}

	return nil
}

// TrustRelationship creates a standard trust relationship for a specific AWS service.
type TrustRelationship struct {
	// Service is the AWS service principal (e.g., "lambda.amazonaws.com").
	Service string

	// Federated is the federated identity provider ARN.
	Federated interface{}

	// AWS is the AWS account or ARN that can assume the role.
	AWS interface{}

	// Conditions are optional conditions for the trust relationship.
	Conditions map[string]interface{}
}

// NewServiceTrustRelationship creates a trust relationship for an AWS service.
func NewServiceTrustRelationship(service string) *TrustRelationship {
	return &TrustRelationship{Service: service}
}

// NewFederatedTrustRelationship creates a trust relationship for a federated identity.
func NewFederatedTrustRelationship(federated interface{}) *TrustRelationship {
	return &TrustRelationship{Federated: federated}
}

// NewAWSTrustRelationship creates a trust relationship for an AWS account or ARN.
func NewAWSTrustRelationship(aws interface{}) *TrustRelationship {
	return &TrustRelationship{AWS: aws}
}

// WithCondition adds a condition to the trust relationship.
func (t *TrustRelationship) WithCondition(conditionType string, condition map[string]interface{}) *TrustRelationship {
	if t.Conditions == nil {
		t.Conditions = make(map[string]interface{})
	}
	t.Conditions[conditionType] = condition
	return t
}

// ToPolicyDocument converts the trust relationship to a PolicyDocument.
func (t *TrustRelationship) ToPolicyDocument() *PolicyDocument {
	principal := make(map[string]interface{})

	if t.Service != "" {
		principal["Service"] = t.Service
	}
	if t.Federated != nil {
		principal["Federated"] = t.Federated
	}
	if t.AWS != nil {
		principal["AWS"] = t.AWS
	}

	statement := &Statement{
		Effect:    EffectAllow,
		Principal: principal,
		Action:    "sts:AssumeRole",
	}

	if len(t.Conditions) > 0 {
		statement.Condition = t.Conditions
	}

	return NewPolicyDocument().AddStatement(statement)
}

// Common AWS service principals for trust relationships.
const (
	ServiceLambda          = "lambda.amazonaws.com"
	ServiceAPIGateway      = "apigateway.amazonaws.com"
	ServiceStepFunctions   = "states.amazonaws.com"
	ServiceEvents          = "events.amazonaws.com"
	ServiceCodeDeploy      = "codedeploy.amazonaws.com"
	ServiceCloudFormation  = "cloudformation.amazonaws.com"
	ServiceEC2             = "ec2.amazonaws.com"
	ServiceECS             = "ecs.amazonaws.com"
	ServiceECSTasksService = "ecs-tasks.amazonaws.com"
	ServiceSNS             = "sns.amazonaws.com"
	ServiceSQS             = "sqs.amazonaws.com"
	ServiceS3              = "s3.amazonaws.com"
	ServiceDynamoDB        = "dynamodb.amazonaws.com"
	ServiceKinesis         = "kinesis.amazonaws.com"
	ServiceFirehose        = "firehose.amazonaws.com"
	ServiceLogs            = "logs.amazonaws.com"
)

// NewLambdaExecutionRole creates a standard Lambda execution role.
func NewLambdaExecutionRole() *Role {
	trustPolicy := NewServiceTrustRelationship(ServiceLambda).ToPolicyDocument()
	return NewRole(trustPolicy)
}

// NewLambdaExecutionRoleWithName creates a named Lambda execution role.
func NewLambdaExecutionRoleWithName(name interface{}) *Role {
	trustPolicy := NewServiceTrustRelationship(ServiceLambda).ToPolicyDocument()
	return NewRoleWithName(name, trustPolicy)
}

// NewStepFunctionsExecutionRole creates a Step Functions execution role.
func NewStepFunctionsExecutionRole() *Role {
	trustPolicy := NewServiceTrustRelationship(ServiceStepFunctions).ToPolicyDocument()
	return NewRole(trustPolicy)
}

// NewAPIGatewayInvocationRole creates an API Gateway invocation role.
func NewAPIGatewayInvocationRole() *Role {
	trustPolicy := NewServiceTrustRelationship(ServiceAPIGateway).ToPolicyDocument()
	return NewRole(trustPolicy)
}

// NewEventsInvocationRole creates an EventBridge/CloudWatch Events invocation role.
func NewEventsInvocationRole() *Role {
	trustPolicy := NewServiceTrustRelationship(ServiceEvents).ToPolicyDocument()
	return NewRole(trustPolicy)
}
