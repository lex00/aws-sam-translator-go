package iam

import (
	"fmt"
)

// PolicyDocumentVersion is the IAM policy document version.
const PolicyDocumentVersion = "2012-10-17"

// Effect constants for policy statements.
const (
	EffectAllow = "Allow"
	EffectDeny  = "Deny"
)

// PolicyDocument represents an IAM policy document following AWS format.
type PolicyDocument struct {
	// Version is the policy language version.
	Version string `json:"Version"`

	// Id is an optional identifier for the policy.
	Id string `json:"Id,omitempty"`

	// Statement is the list of policy statements.
	Statement []*Statement `json:"Statement"`
}

// Statement represents a single statement in an IAM policy document.
type Statement struct {
	// Sid is an optional statement identifier.
	Sid string `json:"Sid,omitempty"`

	// Effect is "Allow" or "Deny".
	Effect string `json:"Effect"`

	// Principal specifies the principal that is allowed or denied access.
	// For trust policies (AssumeRolePolicyDocument).
	Principal interface{} `json:"Principal,omitempty"`

	// NotPrincipal specifies the principal that is NOT allowed or denied access.
	NotPrincipal interface{} `json:"NotPrincipal,omitempty"`

	// Action specifies the actions that are allowed or denied.
	Action interface{} `json:"Action,omitempty"`

	// NotAction specifies the actions that are NOT subject to this statement.
	NotAction interface{} `json:"NotAction,omitempty"`

	// Resource specifies the resources the statement applies to.
	Resource interface{} `json:"Resource,omitempty"`

	// NotResource specifies the resources the statement does NOT apply to.
	NotResource interface{} `json:"NotResource,omitempty"`

	// Condition specifies conditions for when the policy is in effect.
	Condition map[string]interface{} `json:"Condition,omitempty"`
}

// NewPolicyDocument creates a new PolicyDocument with the standard version.
func NewPolicyDocument() *PolicyDocument {
	return &PolicyDocument{
		Version:   PolicyDocumentVersion,
		Statement: make([]*Statement, 0),
	}
}

// NewPolicyDocumentWithId creates a new PolicyDocument with an ID.
func NewPolicyDocumentWithId(id string) *PolicyDocument {
	return &PolicyDocument{
		Version:   PolicyDocumentVersion,
		Id:        id,
		Statement: make([]*Statement, 0),
	}
}

// AddStatement adds a statement to the policy document.
func (d *PolicyDocument) AddStatement(stmt *Statement) *PolicyDocument {
	d.Statement = append(d.Statement, stmt)
	return d
}

// AddStatements adds multiple statements to the policy document.
func (d *PolicyDocument) AddStatements(stmts []*Statement) *PolicyDocument {
	d.Statement = append(d.Statement, stmts...)
	return d
}

// ToMap converts the policy document to a map for CloudFormation.
func (d *PolicyDocument) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"Version": d.Version,
	}

	if d.Id != "" {
		result["Id"] = d.Id
	}

	statements := make([]interface{}, len(d.Statement))
	for i, stmt := range d.Statement {
		statements[i] = stmt.ToMap()
	}
	result["Statement"] = statements

	return result
}

// Validate validates the policy document.
func (d *PolicyDocument) Validate() error {
	if d.Version == "" {
		return fmt.Errorf("version is required in policy document")
	}

	if len(d.Statement) == 0 {
		return fmt.Errorf("at least one Statement is required in policy document")
	}

	for i, stmt := range d.Statement {
		if err := stmt.Validate(); err != nil {
			return fmt.Errorf("invalid statement at index %d: %w", i, err)
		}
	}

	return nil
}

// NewStatement creates a new Statement with the specified effect.
func NewStatement(effect string) *Statement {
	return &Statement{
		Effect: effect,
	}
}

// NewAllowStatement creates a new Allow statement.
func NewAllowStatement() *Statement {
	return NewStatement(EffectAllow)
}

// NewDenyStatement creates a new Deny statement.
func NewDenyStatement() *Statement {
	return NewStatement(EffectDeny)
}

// WithSid sets the statement ID.
func (s *Statement) WithSid(sid string) *Statement {
	s.Sid = sid
	return s
}

// WithPrincipal sets the principal.
func (s *Statement) WithPrincipal(principal interface{}) *Statement {
	s.Principal = principal
	return s
}

// WithServicePrincipal sets a service principal.
func (s *Statement) WithServicePrincipal(service string) *Statement {
	s.Principal = map[string]interface{}{"Service": service}
	return s
}

// WithAWSPrincipal sets an AWS account/ARN principal.
func (s *Statement) WithAWSPrincipal(aws interface{}) *Statement {
	s.Principal = map[string]interface{}{"AWS": aws}
	return s
}

// WithFederatedPrincipal sets a federated principal.
func (s *Statement) WithFederatedPrincipal(federated interface{}) *Statement {
	s.Principal = map[string]interface{}{"Federated": federated}
	return s
}

// WithNotPrincipal sets the NotPrincipal.
func (s *Statement) WithNotPrincipal(notPrincipal interface{}) *Statement {
	s.NotPrincipal = notPrincipal
	return s
}

// WithAction sets the action(s).
func (s *Statement) WithAction(action interface{}) *Statement {
	s.Action = action
	return s
}

// WithActions sets multiple actions.
func (s *Statement) WithActions(actions ...string) *Statement {
	if len(actions) == 1 {
		s.Action = actions[0]
	} else {
		actionList := make([]interface{}, len(actions))
		for i, a := range actions {
			actionList[i] = a
		}
		s.Action = actionList
	}
	return s
}

// WithNotAction sets the NotAction(s).
func (s *Statement) WithNotAction(notAction interface{}) *Statement {
	s.NotAction = notAction
	return s
}

// WithResource sets the resource(s).
func (s *Statement) WithResource(resource interface{}) *Statement {
	s.Resource = resource
	return s
}

// WithResources sets multiple resources.
func (s *Statement) WithResources(resources ...interface{}) *Statement {
	if len(resources) == 1 {
		s.Resource = resources[0]
	} else {
		s.Resource = resources
	}
	return s
}

// WithAllResources sets the resource to "*" (all resources).
func (s *Statement) WithAllResources() *Statement {
	s.Resource = "*"
	return s
}

// WithNotResource sets the NotResource(s).
func (s *Statement) WithNotResource(notResource interface{}) *Statement {
	s.NotResource = notResource
	return s
}

// WithCondition sets a condition.
func (s *Statement) WithCondition(conditionType string, condition map[string]interface{}) *Statement {
	if s.Condition == nil {
		s.Condition = make(map[string]interface{})
	}
	s.Condition[conditionType] = condition
	return s
}

// WithConditions sets all conditions.
func (s *Statement) WithConditions(conditions map[string]interface{}) *Statement {
	s.Condition = conditions
	return s
}

// ToMap converts the statement to a map for CloudFormation.
func (s *Statement) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"Effect": s.Effect,
	}

	if s.Sid != "" {
		result["Sid"] = s.Sid
	}

	if s.Principal != nil {
		result["Principal"] = s.Principal
	}

	if s.NotPrincipal != nil {
		result["NotPrincipal"] = s.NotPrincipal
	}

	if s.Action != nil {
		result["Action"] = s.Action
	}

	if s.NotAction != nil {
		result["NotAction"] = s.NotAction
	}

	if s.Resource != nil {
		result["Resource"] = s.Resource
	}

	if s.NotResource != nil {
		result["NotResource"] = s.NotResource
	}

	if len(s.Condition) > 0 {
		result["Condition"] = s.Condition
	}

	return result
}

// Validate validates the statement.
func (s *Statement) Validate() error {
	if s.Effect != EffectAllow && s.Effect != EffectDeny {
		return fmt.Errorf("effect must be '%s' or '%s'", EffectAllow, EffectDeny)
	}

	// Must have either Action or NotAction
	if s.Action == nil && s.NotAction == nil {
		return fmt.Errorf("statement must have either Action or NotAction")
	}

	// Must have either Resource or NotResource (for resource-based policies)
	// Note: Trust policies (AssumeRolePolicyDocument) don't require Resource
	// So we only validate this if Principal is not set
	if s.Principal == nil && s.NotPrincipal == nil {
		if s.Resource == nil && s.NotResource == nil {
			return fmt.Errorf("statement must have either Resource or NotResource")
		}
	}

	return nil
}

// PolicyDocumentBuilder provides a fluent interface for building policy documents.
type PolicyDocumentBuilder struct {
	document *PolicyDocument
}

// NewPolicyDocumentBuilder creates a new PolicyDocumentBuilder.
func NewPolicyDocumentBuilder() *PolicyDocumentBuilder {
	return &PolicyDocumentBuilder{
		document: NewPolicyDocument(),
	}
}

// WithId sets the policy document ID.
func (b *PolicyDocumentBuilder) WithId(id string) *PolicyDocumentBuilder {
	b.document.Id = id
	return b
}

// AllowActions creates an Allow statement with the specified actions.
func (b *PolicyDocumentBuilder) AllowActions(actions ...string) *StatementBuilder {
	stmt := NewAllowStatement()
	if len(actions) > 0 {
		stmt.WithActions(actions...)
	}
	return &StatementBuilder{
		parent:    b,
		statement: stmt,
	}
}

// DenyActions creates a Deny statement with the specified actions.
func (b *PolicyDocumentBuilder) DenyActions(actions ...string) *StatementBuilder {
	stmt := NewDenyStatement()
	if len(actions) > 0 {
		stmt.WithActions(actions...)
	}
	return &StatementBuilder{
		parent:    b,
		statement: stmt,
	}
}

// AddStatement adds a pre-built statement.
func (b *PolicyDocumentBuilder) AddStatement(stmt *Statement) *PolicyDocumentBuilder {
	b.document.AddStatement(stmt)
	return b
}

// Build returns the constructed PolicyDocument.
func (b *PolicyDocumentBuilder) Build() *PolicyDocument {
	return b.document
}

// StatementBuilder provides a fluent interface for building statements.
type StatementBuilder struct {
	parent    *PolicyDocumentBuilder
	statement *Statement
}

// WithSid sets the statement ID.
func (sb *StatementBuilder) WithSid(sid string) *StatementBuilder {
	sb.statement.WithSid(sid)
	return sb
}

// OnResources sets the resources for this statement.
func (sb *StatementBuilder) OnResources(resources ...interface{}) *StatementBuilder {
	sb.statement.WithResources(resources...)
	return sb
}

// OnAllResources sets the resource to "*".
func (sb *StatementBuilder) OnAllResources() *StatementBuilder {
	sb.statement.WithAllResources()
	return sb
}

// WithCondition adds a condition.
func (sb *StatementBuilder) WithCondition(conditionType string, condition map[string]interface{}) *StatementBuilder {
	sb.statement.WithCondition(conditionType, condition)
	return sb
}

// ForPrincipal sets the principal.
func (sb *StatementBuilder) ForPrincipal(principal interface{}) *StatementBuilder {
	sb.statement.WithPrincipal(principal)
	return sb
}

// ForService sets a service principal.
func (sb *StatementBuilder) ForService(service string) *StatementBuilder {
	sb.statement.WithServicePrincipal(service)
	return sb
}

// ForAWS sets an AWS principal.
func (sb *StatementBuilder) ForAWS(aws interface{}) *StatementBuilder {
	sb.statement.WithAWSPrincipal(aws)
	return sb
}

// Done finishes the statement and returns to the document builder.
func (sb *StatementBuilder) Done() *PolicyDocumentBuilder {
	sb.parent.document.AddStatement(sb.statement)
	return sb.parent
}

// Common policy document helpers.

// NewAssumeRolePolicyForService creates an assume role policy for an AWS service.
func NewAssumeRolePolicyForService(service string) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sts:AssumeRole").
		ForService(service).
		Done().
		Build()
}

// NewAssumeRolePolicyForServices creates an assume role policy for multiple AWS services.
func NewAssumeRolePolicyForServices(services []string) *PolicyDocument {
	var serviceList interface{}
	if len(services) == 1 {
		serviceList = services[0]
	} else {
		list := make([]interface{}, len(services))
		for i, s := range services {
			list[i] = s
		}
		serviceList = list
	}

	stmt := NewAllowStatement().
		WithAction("sts:AssumeRole").
		WithPrincipal(map[string]interface{}{"Service": serviceList})

	return NewPolicyDocument().AddStatement(stmt)
}

// NewAssumeRolePolicyForAccount creates an assume role policy for an AWS account.
func NewAssumeRolePolicyForAccount(accountID string) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sts:AssumeRole").
		ForAWS(fmt.Sprintf("arn:aws:iam::%s:root", accountID)).
		Done().
		Build()
}

// NewAssumeRolePolicyForARN creates an assume role policy for a specific ARN.
func NewAssumeRolePolicyForARN(arn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sts:AssumeRole").
		ForAWS(arn).
		Done().
		Build()
}

// NewResourcePolicy creates a basic resource-based policy.
func NewResourcePolicy(actions []string, resources []interface{}) *PolicyDocument {
	stmt := NewAllowStatement()

	if len(actions) == 1 {
		stmt.Action = actions[0]
	} else {
		actionList := make([]interface{}, len(actions))
		for i, a := range actions {
			actionList[i] = a
		}
		stmt.Action = actionList
	}

	if len(resources) == 1 {
		stmt.Resource = resources[0]
	} else {
		stmt.Resource = resources
	}

	return NewPolicyDocument().AddStatement(stmt)
}

// LambdaBasicExecutionPolicy creates a policy for Lambda basic execution.
func LambdaBasicExecutionPolicy(logGroupArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents").
		OnResources(logGroupArn).
		Done().
		Build()
}

// DynamoDBCrudPolicy creates a policy for DynamoDB CRUD operations.
func DynamoDBCrudPolicy(tableArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions(
			"dynamodb:GetItem",
			"dynamodb:PutItem",
			"dynamodb:UpdateItem",
			"dynamodb:DeleteItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchGetItem",
			"dynamodb:BatchWriteItem",
		).
		OnResources(tableArn).
		Done().
		Build()
}

// S3ReadPolicy creates a policy for S3 read operations.
func S3ReadPolicy(bucketArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("s3:GetObject", "s3:ListBucket").
		OnResources(bucketArn).
		Done().
		Build()
}

// S3CrudPolicy creates a policy for S3 CRUD operations.
func S3CrudPolicy(bucketArn interface{}, objectsArn interface{}) *PolicyDocument {
	stmt1 := NewAllowStatement().
		WithActions("s3:ListBucket", "s3:GetBucketLocation").
		WithResource(bucketArn)

	stmt2 := NewAllowStatement().
		WithActions("s3:GetObject", "s3:PutObject", "s3:DeleteObject").
		WithResource(objectsArn)

	return NewPolicyDocument().AddStatements([]*Statement{stmt1, stmt2})
}

// SQSSendMessagePolicy creates a policy for sending SQS messages.
func SQSSendMessagePolicy(queueArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sqs:SendMessage").
		OnResources(queueArn).
		Done().
		Build()
}

// SQSPollerPolicy creates a policy for polling SQS queues.
func SQSPollerPolicy(queueArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:GetQueueAttributes").
		OnResources(queueArn).
		Done().
		Build()
}

// SNSPublishPolicy creates a policy for publishing to SNS topics.
func SNSPublishPolicy(topicArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("sns:Publish").
		OnResources(topicArn).
		Done().
		Build()
}

// KinesisStreamReadPolicy creates a policy for reading from Kinesis streams.
func KinesisStreamReadPolicy(streamArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions(
			"kinesis:GetRecords",
			"kinesis:GetShardIterator",
			"kinesis:DescribeStream",
			"kinesis:DescribeStreamSummary",
			"kinesis:ListShards",
		).
		OnResources(streamArn).
		Done().
		Build()
}

// StepFunctionsExecutionPolicy creates a policy for Step Functions execution.
func StepFunctionsExecutionPolicy(stateMachineArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("states:StartExecution").
		OnResources(stateMachineArn).
		Done().
		Build()
}

// LambdaInvokePolicy creates a policy for invoking Lambda functions.
func LambdaInvokePolicy(functionArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("lambda:InvokeFunction").
		OnResources(functionArn).
		Done().
		Build()
}

// SecretsManagerReadPolicy creates a policy for reading secrets.
func SecretsManagerReadPolicy(secretArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("secretsmanager:GetSecretValue").
		OnResources(secretArn).
		Done().
		Build()
}

// KMSDecryptPolicy creates a policy for KMS decryption.
func KMSDecryptPolicy(keyArn interface{}) *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions("kms:Decrypt").
		OnResources(keyArn).
		Done().
		Build()
}

// VPCAccessPolicy creates a policy for VPC access (ENI management).
func VPCAccessPolicy() *PolicyDocument {
	return NewPolicyDocumentBuilder().
		AllowActions(
			"ec2:CreateNetworkInterface",
			"ec2:DescribeNetworkInterfaces",
			"ec2:DeleteNetworkInterface",
			"ec2:AssignPrivateIpAddresses",
			"ec2:UnassignPrivateIpAddresses",
		).
		OnAllResources().
		Done().
		Build()
}
