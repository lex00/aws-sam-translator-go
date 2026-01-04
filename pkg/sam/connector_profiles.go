package sam

// ConnectorProfile defines how to generate resources for a source/destination pair.
type ConnectorProfile struct {
	// ResourceType is the CloudFormation resource type to generate.
	// Can be: AWS::IAM::ManagedPolicy, AWS::Lambda::Permission, AWS::SQS::QueuePolicy, AWS::SNS::TopicPolicy
	ResourceType string

	// ReadActions are the IAM actions for Read permission.
	ReadActions []string

	// WriteActions are the IAM actions for Write permission.
	WriteActions []string

	// ReadResourcePatterns are resource patterns for Read permissions.
	// Use %s placeholders for ARN substitution.
	ReadResourcePatterns []ResourcePattern

	// WriteResourcePatterns are resource patterns for Write permissions.
	WriteResourcePatterns []ResourcePattern

	// Principal is the service principal for resource policies.
	Principal string
}

// ResourcePattern defines a resource pattern for IAM policies.
type ResourcePattern struct {
	// Pattern is the resource pattern.
	// Can be "direct" for direct ARN, or a Fn::Sub pattern.
	Pattern string

	// UseArn specifies whether to use the ARN directly or apply pattern.
	UseArn bool

	// SubPattern is the Fn::Sub pattern if not using direct ARN.
	SubPattern string

	// VarName is the variable name in Fn::Sub pattern.
	VarName string
}

// ConnectorProfiles manages all connector profiles.
type ConnectorProfiles struct {
	profiles map[string]map[string]*ConnectorProfile
}

// NewConnectorProfiles creates a new ConnectorProfiles instance.
func NewConnectorProfiles() *ConnectorProfiles {
	p := &ConnectorProfiles{
		profiles: make(map[string]map[string]*ConnectorProfile),
	}
	p.initProfiles()
	return p
}

// GetProfile returns the profile for a source/destination pair.
func (p *ConnectorProfiles) GetProfile(sourceType, destType string) *ConnectorProfile {
	// Normalize serverless types to their CloudFormation equivalents for lookup
	normalizedSource := normalizeResourceType(sourceType)
	normalizedDest := normalizeResourceType(destType)

	if destProfiles, ok := p.profiles[normalizedSource]; ok {
		if profile, ok := destProfiles[normalizedDest]; ok {
			return profile
		}
	}
	return nil
}

// normalizeResourceType normalizes SAM types to CloudFormation types for profile lookup.
func normalizeResourceType(resourceType string) string {
	switch resourceType {
	case TypeServerlessFunction:
		return TypeLambdaFunction
	case TypeServerlessStateMachine:
		return TypeStepFunctionsStateMachine
	case TypeServerlessApi:
		return TypeAPIGatewayRestApi
	case TypeServerlessHttpApi:
		return TypeAPIGatewayV2Api
	default:
		return resourceType
	}
}

// GetActions returns the actions for a permission type.
func (p *ConnectorProfile) GetActions(permission, sourceType, destType string) []string {
	switch permission {
	case "Read":
		return p.ReadActions
	case "Write":
		return p.WriteActions
	default:
		return nil
	}
}

// GetPrincipal returns the service principal for this profile.
func (p *ConnectorProfile) GetPrincipal(sourceType string) string {
	switch sourceType {
	case TypeSNSTopic:
		return "sns.amazonaws.com"
	case TypeS3Bucket:
		return "s3.amazonaws.com"
	case TypeEventsRule:
		return "events.amazonaws.com"
	case TypeAPIGatewayRestApi, TypeServerlessApi:
		return "apigateway.amazonaws.com"
	case TypeAPIGatewayV2Api, TypeServerlessHttpApi:
		return "apigateway.amazonaws.com"
	default:
		return p.Principal
	}
}

// GetResources builds the resource references for IAM policies.
func (p *ConnectorProfile) GetResources(permission string, destArn, sourceArn interface{}, destType, sourceType string) []interface{} {
	var patterns []ResourcePattern
	switch permission {
	case "Read":
		patterns = p.ReadResourcePatterns
	case "Write":
		patterns = p.WriteResourcePatterns
	default:
		return nil
	}

	if len(patterns) == 0 {
		// Default to direct ARN
		return []interface{}{destArn}
	}

	resources := make([]interface{}, 0, len(patterns))
	for _, pattern := range patterns {
		if pattern.UseArn {
			resources = append(resources, destArn)
		} else if pattern.SubPattern != "" {
			// Use Fn::Sub pattern
			resources = append(resources, map[string]interface{}{
				"Fn::Sub": []interface{}{
					pattern.SubPattern,
					map[string]interface{}{
						pattern.VarName: destArn,
					},
				},
			})
		}
	}

	return resources
}

// initProfiles initializes all connector profiles.
func (p *ConnectorProfiles) initProfiles() {
	// Lambda/Function -> DynamoDB Table
	p.addProfile(TypeLambdaFunction, TypeDynamoDBTable, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"dynamodb:GetItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchGetItem",
			"dynamodb:ConditionCheckItem",
			"dynamodb:PartiQLSelect",
		},
		WriteActions: []string{
			"dynamodb:PutItem",
			"dynamodb:UpdateItem",
			"dynamodb:DeleteItem",
			"dynamodb:BatchWriteItem",
			"dynamodb:PartiQLDelete",
			"dynamodb:PartiQLInsert",
			"dynamodb:PartiQLUpdate",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
	})

	// Lambda/Function -> S3 Bucket
	p.addProfile(TypeLambdaFunction, TypeS3Bucket, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"s3:GetObject",
			"s3:GetObjectAcl",
			"s3:GetObjectLegalHold",
			"s3:GetObjectRetention",
			"s3:GetObjectTorrent",
			"s3:GetObjectVersion",
			"s3:GetObjectVersionAcl",
			"s3:GetObjectVersionForReplication",
			"s3:GetObjectVersionTorrent",
			"s3:ListBucket",
			"s3:ListBucketMultipartUploads",
			"s3:ListBucketVersions",
			"s3:ListMultipartUploadParts",
		},
		WriteActions: []string{
			"s3:AbortMultipartUpload",
			"s3:DeleteObject",
			"s3:DeleteObjectVersion",
			"s3:PutObject",
			"s3:PutObjectLegalHold",
			"s3:PutObjectRetention",
			"s3:RestoreObject",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/*", VarName: "DestinationArn"},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/*", VarName: "DestinationArn"},
		},
	})

	// Lambda/Function -> SQS Queue
	p.addProfile(TypeLambdaFunction, TypeSQSQueue, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"sqs:ReceiveMessage",
			"sqs:GetQueueAttributes",
		},
		WriteActions: []string{
			"sqs:DeleteMessage",
			"sqs:SendMessage",
			"sqs:ChangeMessageVisibility",
			"sqs:PurgeQueue",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Lambda/Function -> SNS Topic
	p.addProfile(TypeLambdaFunction, TypeSNSTopic, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"sns:Publish",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Lambda/Function -> Step Functions State Machine
	p.addProfile(TypeLambdaFunction, TypeStepFunctionsStateMachine, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"states:DescribeStateMachine",
			"states:ListExecutions",
		},
		WriteActions: []string{
			"states:StartExecution",
			"states:StartSyncExecution",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Lambda/Function -> Location Place Index
	p.addProfile(TypeLambdaFunction, TypeLocationPlaceIndex, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"geo:SearchPlaceIndexForPosition",
			"geo:SearchPlaceIndexForSuggestions",
			"geo:SearchPlaceIndexForText",
			"geo:GetPlace",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// SNS Topic -> Lambda Function (Push event - requires Lambda permission)
	p.addProfile(TypeSNSTopic, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "sns.amazonaws.com",
	})

	// S3 Bucket -> Lambda Function (Push event - requires Lambda permission)
	p.addProfile(TypeS3Bucket, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "s3.amazonaws.com",
	})

	// SQS Queue -> Lambda Function
	p.addProfile(TypeSQSQueue, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "sqs.amazonaws.com",
	})

	// DynamoDB Table -> Lambda Function (for streams)
	p.addProfile(TypeDynamoDBTable, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"dynamodb:DescribeStream",
			"dynamodb:GetRecords",
			"dynamodb:GetShardIterator",
			"dynamodb:ListStreams",
		},
		ReadResourcePatterns: []ResourcePattern{
			{SubPattern: "${SourceArn}/stream/*", VarName: "SourceArn"},
		},
	})

	// Events Rule -> Lambda Function
	p.addProfile(TypeEventsRule, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "events.amazonaws.com",
	})

	// Events Rule -> SNS Topic
	p.addProfile(TypeEventsRule, TypeSNSTopic, &ConnectorProfile{
		ResourceType: "AWS::SNS::TopicPolicy",
		Principal:    "events.amazonaws.com",
	})

	// Events Rule -> SQS Queue
	p.addProfile(TypeEventsRule, TypeSQSQueue, &ConnectorProfile{
		ResourceType: "AWS::SQS::QueuePolicy",
		Principal:    "events.amazonaws.com",
	})

	// Events Rule -> Step Functions
	p.addProfile(TypeEventsRule, TypeStepFunctionsStateMachine, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"states:StartExecution",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Events Rule -> EventBus
	p.addProfile(TypeEventsRule, TypeEventsEventBus, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"events:PutEvents",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> Lambda Function
	p.addProfile(TypeStepFunctionsStateMachine, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"lambda:InvokeAsync",
			"lambda:InvokeFunction",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> Step Functions
	p.addProfile(TypeStepFunctionsStateMachine, TypeStepFunctionsStateMachine, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"states:StartExecution",
			"states:StartSyncExecution",
		},
		ReadActions: []string{
			"states:DescribeExecution",
			"states:StopExecution",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// API Gateway -> Lambda Function
	p.addProfile(TypeAPIGatewayRestApi, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "apigateway.amazonaws.com",
	})

	// API Gateway V2 -> Lambda Function
	p.addProfile(TypeAPIGatewayV2Api, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::Lambda::Permission",
		Principal:    "apigateway.amazonaws.com",
	})

	// SNS -> SQS
	p.addProfile(TypeSNSTopic, TypeSQSQueue, &ConnectorProfile{
		ResourceType: "AWS::SQS::QueuePolicy",
		Principal:    "sns.amazonaws.com",
	})

	// AppSync GraphQL API -> Lambda
	p.addProfile(TypeAppSyncGraphQLApi, TypeLambdaFunction, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"lambda:InvokeFunction",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// AppSync GraphQL API -> DynamoDB Table
	p.addProfile(TypeAppSyncGraphQLApi, TypeDynamoDBTable, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"dynamodb:GetItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchGetItem",
		},
		WriteActions: []string{
			"dynamodb:PutItem",
			"dynamodb:UpdateItem",
			"dynamodb:DeleteItem",
			"dynamodb:BatchWriteItem",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
	})

	// AppSync GraphQL API -> EventBus
	p.addProfile(TypeAppSyncGraphQLApi, TypeEventsEventBus, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"events:PutEvents",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Lambda -> EventBus
	p.addProfile(TypeLambdaFunction, TypeEventsEventBus, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"events:PutEvents",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> DynamoDB Table
	p.addProfile(TypeStepFunctionsStateMachine, TypeDynamoDBTable, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"dynamodb:GetItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchGetItem",
			"dynamodb:ConditionCheckItem",
		},
		WriteActions: []string{
			"dynamodb:PutItem",
			"dynamodb:UpdateItem",
			"dynamodb:DeleteItem",
			"dynamodb:BatchWriteItem",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/index/*", VarName: "DestinationArn"},
		},
	})

	// Step Functions -> SQS
	p.addProfile(TypeStepFunctionsStateMachine, TypeSQSQueue, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"sqs:SendMessage",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> SNS
	p.addProfile(TypeStepFunctionsStateMachine, TypeSNSTopic, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"sns:Publish",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> EventBus
	p.addProfile(TypeStepFunctionsStateMachine, TypeEventsEventBus, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		WriteActions: []string{
			"events:PutEvents",
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
		},
	})

	// Step Functions -> S3
	p.addProfile(TypeStepFunctionsStateMachine, TypeS3Bucket, &ConnectorProfile{
		ResourceType: "AWS::IAM::ManagedPolicy",
		ReadActions: []string{
			"s3:GetObject",
			"s3:ListBucket",
		},
		WriteActions: []string{
			"s3:PutObject",
			"s3:DeleteObject",
		},
		ReadResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/*", VarName: "DestinationArn"},
		},
		WriteResourcePatterns: []ResourcePattern{
			{UseArn: true},
			{SubPattern: "${DestinationArn}/*", VarName: "DestinationArn"},
		},
	})
}

// addProfile adds a profile for a source/destination pair.
func (p *ConnectorProfiles) addProfile(sourceType, destType string, profile *ConnectorProfile) {
	if _, ok := p.profiles[sourceType]; !ok {
		p.profiles[sourceType] = make(map[string]*ConnectorProfile)
	}
	p.profiles[sourceType][destType] = profile
}
