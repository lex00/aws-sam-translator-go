# Architecture Overview

This document describes the architecture of aws-sam-translator-go, a Go implementation of the AWS SAM to CloudFormation translator.

## High-Level Flow

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  SAM Template   │────>│     Parser      │────>│  types.Template │
│  (YAML/JSON)    │     │  (pkg/parser)   │     │                 │
└─────────────────┘     └─────────────────┘     └────────┬────────┘
                                                         │
                        ┌─────────────────┐              │
                        │  Plugin System  │<─────────────┤
                        │  BeforeTransform│              │
                        └────────┬────────┘              │
                                 │                       │
                        ┌────────▼────────┐              │
                        │   Intrinsics    │<─────────────┘
                        │    Resolver     │
                        │ (pkg/intrinsics)│
                        └────────┬────────┘
                                 │
                        ┌────────▼────────┐
                        │ SAM Transformers│
                        │   (pkg/sam)     │
                        │  - Function     │
                        │  - StateMachine │
                        │  - Connector    │
                        │  - SimpleTable  │
                        │  - LayerVersion │
                        └────────┬────────┘
                                 │
                        ┌────────▼────────┐
                        │  Plugin System  │
                        │ AfterTransform  │
                        └────────┬────────┘
                                 │
                        ┌────────▼────────┐
                        │  CloudFormation │
                        │    Template     │
                        └─────────────────┘
```

## Package Overview

### Core Packages

| Package | Purpose | Key Types |
|---------|---------|-----------|
| `pkg/types` | Core data structures | `Template`, `Resource`, `Parameter`, `Output` |
| `pkg/parser` | Template parsing | `Parse()`, YAML/JSON handlers |
| `pkg/intrinsics` | Intrinsic function resolution | `Registry`, `ResolveContext`, `Action` |
| `pkg/sam` | SAM resource transformers | `FunctionTransformer`, `ConnectorTransformer`, etc. |
| `pkg/plugins` | Plugin system | `Plugin`, `Registry`, built-in plugins |
| `pkg/openapi` | OpenAPI/Swagger generation | `Generator`, `Route`, `GenerateSwagger()`, `GenerateOpenAPI3()` |
| `pkg/translator` | ID and ARN generation | `LogicalIdGenerator`, `ArnGenerator` |
| `pkg/policy` | Policy template expansion | `Processor`, 81 embedded templates |
| `pkg/region` | AWS region/partition config | `GetPartitionForRegion()`, `GetArnPartition()` |

### Model Packages

| Package | Purpose |
|---------|---------|
| `pkg/model/iam` | IAM Role, Policy, ManagedPolicy models |
| `pkg/model/lambda` | Lambda Function, Version, Alias, Permission, EventSourceMapping |
| `pkg/model/eventsources/push` | Push event sources (S3, SNS, API, etc.) |
| `pkg/model/eventsources/pull` | Pull event sources (SQS, Kinesis, DynamoDB, etc.) |
| `pkg/cloudformation/*` | CloudFormation resource implementations |

## Component Details

### 1. Parser (`pkg/parser`)

The parser handles YAML and JSON SAM templates with full CloudFormation intrinsic function support.

```go
// Parse auto-detects format and returns a Template
template, err := parser.Parse([]byte(yamlContent))
```

**Key Features:**
- Auto-detection of YAML vs JSON format
- Custom YAML tag handlers for 17 intrinsic short-form tags (!Ref, !Sub, !GetAtt, etc.)
- Source location tracking (line/column) for error reporting
- Automatic conversion of `!GetAtt Resource.Attr` to array form

**Supported Intrinsic Tags:**
- `!Ref`, `!Sub`, `!GetAtt`, `!Join`, `!If`, `!Select`
- `!FindInMap`, `!Base64`, `!Cidr`, `!GetAZs`
- `!ImportValue`, `!Split`, `!Transform`
- `!And`, `!Equals`, `!Not`, `!Or`, `!Condition`

### 2. Intrinsics Resolver (`pkg/intrinsics`)

Resolves CloudFormation intrinsic functions using pre-order tree traversal.

```go
ctx := intrinsics.NewResolveContext(template)
ctx.SetParameter("Environment", "prod")

registry := intrinsics.NewRegistry()
result, err := registry.Resolve(ctx, value)
```

**Architecture:**
- **Action Interface**: Each intrinsic (Ref, Fn::Sub, etc.) implements `Action`
- **Registry**: Manages action handlers, executes resolution
- **ResolveContext**: Holds parameters, resources, pseudo-parameters
- **Dependency Tracking**: Builds resource dependency graphs

**Resolved Intrinsics:**
- `Ref` - Parameters, resources, pseudo-parameters
- `Fn::Sub` - String/array forms with variable substitution
- `Fn::GetAtt` - Resource attribute lookups
- `Fn::FindInMap` - Mapping resolution
- `Fn::Join` - Static evaluation when all values are strings

**Pass-Through Intrinsics** (preserved for CloudFormation runtime):
- `Fn::ImportValue`, `Fn::GetAZs`, `Fn::If`, `Fn::Select`, `Fn::Split`, `Fn::Base64`

### 3. SAM Transformers (`pkg/sam`)

Each SAM resource type has a dedicated transformer that converts it to CloudFormation resources.

#### Function Transformer

Transforms `AWS::Serverless::Function` → multiple CloudFormation resources:

```
AWS::Serverless::Function
    │
    ├──> AWS::Lambda::Function
    ├──> AWS::IAM::Role (if Role not specified)
    ├──> AWS::Lambda::Version (if AutoPublishAlias)
    ├──> AWS::Lambda::Alias (if AutoPublishAlias)
    ├──> AWS::Lambda::Permission (for each push event)
    ├──> AWS::Lambda::EventSourceMapping (for each pull event)
    └──> Event-specific resources (API Gateway, S3, SNS, etc.)
```

#### StateMachine Transformer

Transforms `AWS::Serverless::StateMachine`:

```
AWS::Serverless::StateMachine
    │
    ├──> AWS::StepFunctions::StateMachine
    ├──> AWS::IAM::Role (if Role not specified)
    ├──> AWS::Logs::LogGroup (if logging enabled)
    └──> Event-specific resources (Schedule, CloudWatch, etc.)
```

#### Connector Transformer

Transforms `AWS::Serverless::Connector` into IAM policies:

```
AWS::Serverless::Connector
    │
    ├──> AWS::IAM::Policy (attached to source role)
    ├──> AWS::Lambda::Permission (for invocation permissions)
    ├──> AWS::SQS::QueuePolicy (for SQS targets)
    └──> AWS::SNS::TopicPolicy (for SNS targets)
```

### 4. Plugin System (`pkg/plugins`)

Plugins modify templates before and after SAM transformation.

```go
type Plugin interface {
    Name() string
    Priority() int  // Lower = runs first
    BeforeTransform(template *types.Template) error
    AfterTransform(template *types.Template) error
}
```

**Built-in Plugins:**

| Plugin | Priority | Purpose |
|--------|----------|---------|
| GlobalsPlugin | 100 | Merges Globals section into resources |
| ImplicitRestApiPlugin | 300 | Creates implicit REST API from Api events |
| ImplicitHttpApiPlugin | 310 | Creates implicit HTTP API from HttpApi events |
| PolicyTemplatesPlugin | 400 | Expands SAM policy templates |
| DefaultDefinitionBodyPlugin | 500 | Generates OpenAPI specs from function events |

### 5. Event Source Handlers

#### Push Events (`pkg/model/eventsources/push`)

Push events create Lambda permissions and event source configurations:

| Event | Creates |
|-------|---------|
| S3 | Lambda::Permission, S3 bucket notification config |
| SNS | Lambda::Permission, SNS::Subscription |
| Api | Lambda::Permission, API Gateway resources |
| HttpApi | Lambda::Permission, API Gateway V2 resources |
| Schedule | Lambda::Permission, Events::Rule |
| CloudWatch | Lambda::Permission, Events::Rule |
| Cognito | Lambda::Permission, Cognito trigger config |
| IoT | Lambda::Permission, IoT::TopicRule |

#### Pull Events (`pkg/model/eventsources/pull`)

Pull events create EventSourceMapping resources:

| Event | Source |
|-------|--------|
| SQS | SQS Queue |
| Kinesis | Kinesis Stream |
| DynamoDB | DynamoDB Stream |
| DocumentDB | DocumentDB change stream |
| MSK | Managed Kafka cluster |
| MQ | Amazon MQ broker |
| SelfManagedKafka | Self-hosted Kafka |
| ScheduleV2 | EventBridge Scheduler |

### 6. Policy Template Processor (`pkg/policy`)

Expands SAM policy templates with parameter substitution.

```go
processor, _ := policy.New()

// Expand template with parameters
definition, _ := processor.Expand("DynamoDBCrudPolicy", map[string]interface{}{
    "TableName": "MyTable",
})
```

**Features:**
- 81 embedded policy templates (DynamoDB, S3, Lambda, SQS, SNS, etc.)
- Supports `Ref` and `Fn::Sub` in template parameters
- Parameter validation with required field checking
- Returns complete policy document or just statements

### 7. ID and ARN Generators (`pkg/translator`)

#### Logical ID Generator

Generates stable, deterministic CloudFormation logical IDs:

```go
gen := NewLogicalIdGenerator()
id := gen.Generate("MyFunction", "Role")  // "MyFunctionRole"
```

- Uses SHA256 hashing for long names
- Ensures uniqueness within template
- Maintains stability across transformations

#### ARN Generator

Generates partition-aware ARNs:

```go
gen := NewArnGenerator("us-east-1")
arn := gen.Lambda("123456789012", "my-function")
// arn:aws:lambda:us-east-1:123456789012:function:my-function

gen := NewArnGenerator("cn-north-1")
arn := gen.Lambda("123456789012", "my-function")
// arn:aws-cn:lambda:cn-north-1:123456789012:function:my-function
```

**Supported Services:**
Lambda, API Gateway, IAM, S3, DynamoDB, SNS, SQS, Kinesis, Step Functions, EventBridge, CloudWatch, Secrets Manager, KMS, Cognito, CodeDeploy

### 8. Region Configuration (`pkg/region`)

Handles AWS partition detection and region utilities.

```go
partition := region.GetPartitionForRegion("us-east-1")  // "aws"
partition := region.GetPartitionForRegion("cn-north-1") // "aws-cn"
partition := region.GetPartitionForRegion("us-gov-west-1") // "aws-us-gov"

suffix := region.GetDNSSuffix("aws")     // "amazonaws.com"
suffix := region.GetDNSSuffix("aws-cn")  // "amazonaws.com.cn"
```

## Design Patterns

### Strategy Pattern
Intrinsic action handlers implement the `Action` interface, allowing different resolution strategies.

### Registry Pattern
Plugin and action registries manage extensibility, enabling runtime registration of handlers.

### Builder Pattern
CloudFormation models use builder methods for fluent configuration.

### Visitor Pattern
The intrinsic resolver traverses the template tree, visiting and resolving nodes.

### Factory Pattern
Transformers are created via factory functions (`NewFunctionTransformer`, etc.).

## Dependency Flow

```
types ─────────────> parser
  │                    │
  │                    ▼
  │               intrinsics ◄──── region
  │                    │
  │                    ▼
  │                 plugins
  │                    │
  │                    ▼
  ├───────────────> sam ◄───────── policy
  │                    │
  │                    ▼
  │               translator
  │                    │
  │                    ▼
  └──────────────> model/* ◄──── cloudformation/*
```

## Extension Points

1. **New SAM Resource Types**: Implement a transformer in `pkg/sam/`
2. **New Event Sources**: Add handlers in `pkg/model/eventsources/`
3. **New Intrinsic Actions**: Implement `Action` interface, register in `Registry`
4. **New Plugins**: Implement `Plugin` interface, register in plugin registry
5. **New Policy Templates**: Add to `pkg/policy/templates.json`

## Error Handling

Custom error types in `pkg/errors/`:

- `InvalidDocumentException` - Template structure errors
- `InvalidResourceException` - Resource configuration errors
- `InvalidEventException` - Event source configuration errors

All errors include context (resource name, property path) for debugging.
