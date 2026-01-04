# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **AWS::Serverless::Api Transformer** (Phase 6A - #15)
  - `pkg/sam/api.go` - Complete Api transformer implementation
  - `pkg/sam/api_test.go` - 27 comprehensive tests
  - RestApi with Swagger/OpenAPI 2.0 and 3.0 definition support
  - DefinitionBody and DefinitionUri (S3 location) handling
  - Stage configuration with variables, caching, tracing
  - Deployment management with hash-based stable IDs
  - Authorizers (Cognito User Pools, Lambda TOKEN/REQUEST)
  - Endpoint configuration (EDGE, REGIONAL, PRIVATE)
  - CORS configuration support
  - Binary media types
  - Access logging configuration
  - Method settings (logging, metrics, caching, throttling)
  - Gateway responses
  - MinimumCompressionSize, FailOnWarnings, DisableExecuteApiEndpoint
  - Tags support

- **AWS::Serverless::Function Transformer** (Phase 5A - #49)
  - `pkg/sam/function.go` - Complete Function transformer implementation (1,440 lines)
  - `pkg/sam/function_test.go` - Comprehensive tests
  - Event source extraction and transformation for all 18 event types
  - IAM role auto-creation with AssumeRolePolicyDocument
  - Policy attachment from Policies property (SAM policy templates, managed ARNs, inline policies)
  - Lambda alias and version management
  - Deployment preferences with CodeDeploy integration
  - Provisioned concurrency configuration
  - VPC configuration support
  - Environment variables and tags
  - Dead letter queue configuration
  - Tracing and logging settings

- **AWS::Serverless::StateMachine Transformer** (Phase 6C - #48)
  - `pkg/sam/statemachine.go` - Complete StateMachine transformer implementation (631 lines)
  - `pkg/sam/statemachine_test.go` - Comprehensive tests
  - Definition and DefinitionUri resolution
  - DefinitionSubstitutions for variable replacement
  - IAM role auto-creation with Step Functions assume role policy
  - Policies transformation (SAM policy templates, managed ARNs, inline)
  - Logging configuration with CloudWatch Logs integration
  - X-Ray tracing configuration
  - Tags propagation
  - Event source support (Schedule, CloudWatchEvent, EventBridgeRule, Api)

- **AWS::Serverless::SimpleTable Transformer** (Phase 6A - #40)
  - `pkg/sam/simpletable.go` - SimpleTable to DynamoDB::Table transformer
  - `pkg/sam/simpletable_test.go` - Comprehensive tests
  - Primary key configuration with AttributeName and Type
  - Provisioned throughput settings
  - SSE specification support
  - Table name and tags

- **AWS::Serverless::LayerVersion Transformer** (Phase 6B - #40)
  - `pkg/sam/layerversion.go` - LayerVersion transformer
  - `pkg/sam/layerversion_test.go` - Comprehensive tests
  - ContentUri resolution (S3 location or local path)
  - Compatible runtimes and architectures
  - License info and description
  - Layer name and retention policy

- **Plugin System** (Phase 8 - #39, #20)
  - `pkg/plugins/plugin.go` - Plugin interface and Registry with priority-based execution
  - `pkg/plugins/registry_test.go` - Registry tests
  - `pkg/plugins/globals.go` - GlobalsPlugin: Merges Globals section properties into resources
  - `pkg/plugins/globals_test.go` - Globals tests
  - `pkg/plugins/implicit_api.go` - ImplicitRestApiPlugin: Creates implicit REST API from function Api events
  - `pkg/plugins/implicit_api_test.go` - Implicit REST API tests
  - `pkg/plugins/implicit_httpapi.go` - ImplicitHttpApiPlugin: Creates implicit HTTP API from function HttpApi events
  - `pkg/plugins/implicit_httpapi_test.go` - Implicit HTTP API tests
  - `pkg/plugins/policy_templates.go` - PolicyTemplatesPlugin: Expands SAM policy templates in Policies property
  - `pkg/plugins/policy_templates_test.go` - Policy templates tests
  - `pkg/plugins/default_definition_body.go` - DefaultDefinitionBodyPlugin: Sets default OpenAPI definition body
  - `pkg/plugins/default_definition_body_test.go` - Default definition body tests
  - BeforeTransform and AfterTransform hooks for template modification

- **Push Event Source Handlers** (Phase 4A)
  - `pkg/model/eventsources/push/s3.go` - S3 event source with bucket notifications
  - `pkg/model/eventsources/push/sns.go` - SNS event source with topic subscriptions
  - `pkg/model/eventsources/push/api.go` - REST API Gateway event source
  - `pkg/model/eventsources/push/httpapi.go` - HTTP API Gateway (v2) event source
  - `pkg/model/eventsources/push/schedule.go` - EventBridge Schedule event source
  - `pkg/model/eventsources/push/cloudwatch.go` - CloudWatch Events/EventBridge Rules
  - `pkg/model/eventsources/push/cognito.go` - Cognito User Pool triggers
  - `pkg/model/eventsources/push/iot.go` - IoT Rule event source
  - Full test coverage for all push event handlers

- **Pull Event Source Handlers** (Phase 4B)
  - `pkg/model/eventsources/pull/sqs.go` - SQS queue polling with batch settings
  - `pkg/model/eventsources/pull/kinesis.go` - Kinesis stream with starting position, parallelization
  - `pkg/model/eventsources/pull/dynamodb.go` - DynamoDB Streams with batch and filter settings
  - `pkg/model/eventsources/pull/documentdb.go` - DocumentDB change streams
  - `pkg/model/eventsources/pull/msk.go` - Amazon MSK (Managed Streaming for Kafka)
  - `pkg/model/eventsources/pull/mq.go` - Amazon MQ (ActiveMQ, RabbitMQ)
  - `pkg/model/eventsources/pull/cloudwatchlogs.go` - CloudWatch Logs subscription
  - `pkg/model/eventsources/pull/selfmanagedkafka.go` - Self-managed Kafka clusters
  - `pkg/model/eventsources/pull/schedulev2.go` - EventBridge Scheduler (v2)
  - Full test coverage for all pull event handlers

- **AWS::Serverless::Connector Transformer** (Phase 7A - #18)
  - `pkg/sam/connector.go` - Complete Connector transformer implementation
  - `pkg/sam/connector_profiles.go` - Connector profiles for all service pairs
  - `pkg/sam/connector_test.go` - 23 comprehensive tests for connector functionality
  - Source/Destination resource mapping with automatic type resolution
  - Permission type resolution (Read, Write) with profile-based actions
  - IAM policy generation for all supported service pairs:
    - Lambda/Function -> DynamoDB, S3, SQS, SNS, Step Functions, Location, EventBus
    - SNS/S3/SQS/Events Rule -> Lambda (Lambda permissions)
    - Events Rule -> SQS (Queue policies), SNS (Topic policies), Step Functions, EventBus
    - Step Functions -> Lambda, DynamoDB, SQS, SNS, S3, EventBus, Step Functions
    - API Gateway/HTTP API -> Lambda
    - AppSync GraphQL API -> Lambda, DynamoDB, EventBus
  - Embedded connector extraction from resource Connectors property
  - Policy consolidation for multiple permissions on same resource pair
  - SAM type normalization (Serverless::Function -> Lambda::Function, etc.)
  - Proper role reference extraction for policy attachment

- **CloudFormation Model Tests for IoT and Cognito** (Phase 4A - #46)
  - `pkg/cloudformation/iot/topic_rule_test.go` - 100% test coverage for IoT TopicRule model
  - `pkg/cloudformation/cognito/user_pool_test.go` - 100% test coverage for Cognito LambdaConfig and trigger types
  - Tests for NewTopicRule constructor, builder methods, and ToCloudFormation output
  - Tests for all 10 Cognito trigger type constants and GetLambdaConfigProperty function

- **Logical ID and ARN Generators** (Phase 2B-2C - #5)
  - `pkg/translator/logical_id.go` - Hash-based LogicalIdGenerator with SHA256 for deterministic, stable CloudFormation logical IDs
  - `pkg/translator/arn.go` - Partition-aware ArnGenerator supporting aws, aws-cn, aws-us-gov partitions
  - `pkg/translator/verify_id.go` - IDVerifier for logical ID validation, duplicate detection, and ARN verification
  - Support for 15+ AWS services: Lambda, API Gateway, IAM, S3, DynamoDB, SNS, SQS, Kinesis, Step Functions, EventBridge, CloudWatch, Secrets Manager, KMS, Cognito, CodeDeploy
  - API deployment ID recalculation on OpenAPI spec change
  - IDStabilityChecker for deterministic ID generation verification

- **Intrinsics Resolver with Tree Traversal** (Phase 2A - #4)
  - `pkg/intrinsics/resolver.go` - Pre-order tree traversal resolver for nested intrinsic function resolution
  - `pkg/intrinsics/resource_refs.go` - DependencyTracker for building resource dependency graphs
  - Logical ID mutation support for SAM transformations
  - Placeholder protection for CloudFormation runtime intrinsics (Fn::ImportValue, Fn::GetAZs)
  - Topological sort for processing resources in dependency order
  - ResourceRefCollector for scanning templates for Ref and GetAtt references

- **YAML/JSON Parser with Intrinsic Support** (Phase 1B - #2)
  - `pkg/parser/yaml.go` - Custom YAML tag handling for all CloudFormation intrinsic short-form tags (!Ref, !Sub, !GetAtt, !Join, !If, !Select, !FindInMap, !Base64, !Cidr, !GetAZs, !ImportValue, !Split, !Transform, !And, !Equals, !Not, !Or, !Condition)
  - `pkg/parser/json.go` - JSON parsing with intrinsic structure preservation
  - `pkg/parser/intrinsics.go` - Intrinsic function detection, validation, and structure checking
  - `pkg/parser/parser.go` - Enhanced parser with source location tracking (line/column) for error reporting
  - Automatic conversion of !GetAtt dot notation (`!GetAtt Resource.Attr`) to array form
  - Template validation and auto-detection of YAML vs JSON format
  - 24 comprehensive tests

- **Core Types and Region Configuration** (Phase 1A - #1)
  - `pkg/region/config.go` - AWS partition detection for `aws`, `aws-cn`, and `aws-us-gov` partitions
  - Region-to-partition mapping functions (`GetPartitionForRegion`, `GetArnPartition`, `GetDNSSuffix`)
  - Region validation and default region handling
  - `pkg/utils/utils.go` - Extended with sorting functions (`SortedKeys`, `SortedResourceKeys`, `SortStringSlice`, `UniqueStrings`) for deterministic output
  - Comparison utilities (`DeepEqual`, `MapContains`, `StringSliceContains`, `MergeMaps`)
  - Safe accessor functions (`GetStringValue`, `GetMapValue`, `GetSliceValue`)
  - Comprehensive test coverage for errors, types, utils, and region packages

- **Test Fixtures** (Phase 10 - #23)
  - Ported 2,583 test fixtures from upstream Python [aws-sam-translator](https://github.com/aws/serverless-application-model)
  - 775 input SAM YAML templates (513 success + 262 error cases)
  - 778 expected CloudFormation JSON outputs (default partition)
  - 515 partition-specific outputs for `aws-cn`
  - 515 partition-specific outputs for `aws-us-gov`
  - Coverage for all SAM resource types: Function, Api, HttpApi, SimpleTable, LayerVersion, StateMachine, Connector, Application, GraphQLApi

- **Intrinsic Function Handlers** (Phase 1C - #3)
  - `pkg/intrinsics/actions.go` - Core types: `Action` interface, `Registry`, `ResolveContext`
  - `pkg/intrinsics/ref.go` - `RefAction` for `Ref` intrinsics (parameters, resources, pseudo-parameters)
  - `pkg/intrinsics/sub.go` - `SubAction` for `Fn::Sub` variable substitution (string and array forms)
  - `pkg/intrinsics/getatt.go` - `GetAttAction` for `Fn::GetAtt` resource attribute lookups
  - `pkg/intrinsics/findinmap.go` - `FindInMapAction` for `Fn::FindInMap` mapping resolution
  - `pkg/intrinsics/passthrough.go` - Pass-through handlers for `Fn::Join`, `Fn::If`, `Fn::Select`, `Fn::Base64`, `Fn::GetAZs`, `Fn::Split`, `Fn::ImportValue`, `Condition`
  - `ResolveContextOptions` for configurable pseudo-parameters (AccountId, Region, StackName, etc.)
  - `NoValue` sentinel type and `IsNoValue()` helper for `AWS::NoValue` handling
  - Static evaluation for `Fn::Join` when all values are strings
  - Nested intrinsic function resolution
  - 113 comprehensive tests

- **Policy Template Processor** (Phase 2D - #6)
  - `pkg/policy/processor.go` - Core processor with template loading, parameter substitution, and IAM statement generation
  - `pkg/policy/templates.json` - 81 SAM policy templates embedded at compile time via `go:embed`
  - Support for `Ref` and `Fn::Sub` intrinsic function parameter substitution
  - API: `New()`, `Expand()`, `ExpandStatements()`, `HasTemplate()`, `GetTemplate()`, `TemplateNames()`

- **IAM CloudFormation Resource Models** (Phase 3A - #7)
  - `pkg/model/iam/role.go` - IAM Role with AssumeRolePolicyDocument, ManagedPolicyArns, Policies
  - `pkg/model/iam/policy.go` - IAM Policy with PolicyDocument
  - `pkg/model/iam/managed.go` - IAM ManagedPolicy
  - `pkg/model/iam/document.go` - PolicyDocument structure with Statement support

- **Lambda CloudFormation Resource Models** (Phase 3B - #8)
  - `pkg/model/lambda/function.go` - Lambda Function with all configuration options
  - `pkg/model/lambda/version.go` - Lambda Version for function versioning
  - `pkg/model/lambda/alias.go` - Lambda Alias with routing configuration
  - `pkg/model/lambda/permission.go` - Lambda Permission for resource-based policies
  - `pkg/model/lambda/eventsourcemapping.go` - EventSourceMapping for event source triggers
  - `pkg/model/lambda/layer.go` - Lambda LayerVersion

- **API Gateway CloudFormation Resource Models** (Phase 3C - #9)
  - `pkg/cloudformation/apigateway/` - REST API resources:
    - `restapi.go` - RestApi with endpoint configuration, CORS, and policies
    - `stage.go` - Stage with caching, logging, and canary settings
    - `deployment.go` - Deployment with stage descriptions
    - `authorizer.go` - TOKEN, REQUEST, and COGNITO authorizer types
    - `method.go` - HTTP method with integration configuration
    - `resource.go` - API resource path definitions
  - `pkg/cloudformation/apigatewayv2/` - HTTP/WebSocket API resources:
    - `api.go` - HTTP/WebSocket API with CORS and protocol settings
    - `stage.go` - Stage with auto-deploy and access logging
    - `integration.go` - Lambda, HTTP, and AWS service integrations
    - `route.go` - Route definitions with authorization
    - `authorizer.go` - JWT and REQUEST authorizer types

- **Additional CloudFormation Resource Models** (Phase 3D - #10)
  - `pkg/cloudformation/dynamodb/table.go` - DynamoDB Table with GSI, LSI, streams, TTL
  - `pkg/cloudformation/events/rule.go` - EventBridge Rule with targets
  - `pkg/cloudformation/stepfunctions/statemachine.go` - Step Functions StateMachine with logging
  - `pkg/cloudformation/sns/topic.go` - SNS Topic
  - `pkg/cloudformation/sns/subscription.go` - SNS Subscription
  - `pkg/cloudformation/sqs/queue.go` - SQS Queue with FIFO and DLQ support
  - `pkg/cloudformation/logs/loggroup.go` - CloudWatch Logs LogGroup
  - `pkg/cloudformation/s3/bucket.go` - S3 Bucket with notification configuration

- **Core Infrastructure**
  - `pkg/types/` - SAM and CloudFormation type definitions
  - `pkg/errors/` - Error types (InvalidDocumentException, InvalidResourceException, InvalidEventException)
  - `pkg/parser/` - Template parser foundation
  - `pkg/intrinsics/` - Intrinsic function handling foundation
  - `pkg/translator/` - Main translator foundation

### Policy Templates

81 policy templates ported from [aws-sam-translator](https://github.com/aws/serverless-application-model), including:

- **DynamoDB**: CrudPolicy, ReadPolicy, WritePolicy, StreamReadPolicy
- **S3**: CrudPolicy, ReadPolicy, WritePolicy, FullAccessPolicy
- **Lambda**: InvokePolicy
- **SQS**: PollerPolicy, SendMessagePolicy
- **SNS**: CrudPolicy, PublishMessagePolicy
- **Kinesis**: CrudPolicy, StreamReadPolicy
- **KMS**: DecryptPolicy, EncryptPolicy
- **Secrets Manager**: GetSecretValuePolicy, RotationPolicy
- **Step Functions**: ExecutionPolicy
- And 60+ more...
