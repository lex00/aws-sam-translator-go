# aws-sam-translator-go

A Go port of [aws-sam-translator](https://github.com/aws/serverless-application-model) for transforming AWS SAM templates to CloudFormation.

## Status

**Work in Progress** - See the [Implementation Roadmap](https://github.com/lex00/aws-sam-translator-go/issues/24) for current progress.

### Completed

- [x] **Phase 1A**: Core types, region/partition config, utility functions
- [x] **Phase 1B**: YAML/JSON parser with CloudFormation intrinsic tag support
- [x] **Phase 1C**: Intrinsic function handlers (Ref, Fn::Sub, Fn::GetAtt, Fn::FindInMap, pass-through handlers)
- [x] **Phase 2A**: Intrinsics resolver with tree traversal, dependency tracking
- [x] **Phase 2B-2C**: Logical ID and ARN generators with partition support
- [x] **Phase 2D**: Policy template processor with 81 SAM policy templates
- [x] **Phase 3A**: IAM CloudFormation resource models (Role, Policy, ManagedPolicy)
- [x] **Phase 3B**: Lambda CloudFormation resource models (Function, Version, Alias, Permission, EventSourceMapping, LayerVersion)
- [x] **Phase 3C**: API Gateway CloudFormation resource models (RestApi, Stage, Deployment, Authorizer, Method, Resource; V2 Api, Stage, Integration, Route, Authorizer)
- [x] **Phase 3D**: Additional CloudFormation resource models (DynamoDB, EventBridge, Step Functions, SNS, SQS, S3, CloudWatch Logs)
- [x] **Phase 4A-4B**: Event source handlers - 9 push sources (S3, SNS, API, HttpApi, Schedule, CloudWatch, Cognito, IoT) and 9 pull sources (SQS, Kinesis, DynamoDB, DocumentDB, MSK, MQ, CloudWatchLogs, SelfManagedKafka, ScheduleV2)
- [x] **Phase 5A**: AWS::Serverless::Function transformer with event sources, IAM roles, aliases, deployment preferences
- [x] **Phase 5B-5C**: AWS::Serverless::SimpleTable and AWS::Serverless::LayerVersion transformers
- [x] **Phase 6A**: AWS::Serverless::Api transformer with Swagger/OpenAPI, authorizers, CORS, caching
- [x] **Phase 6B**: AWS::Serverless::HttpApi transformer with JWT/Lambda authorizers, CORS, custom domains
- [x] **Phase 6C**: AWS::Serverless::StateMachine transformer with logging, tracing, policies
- [x] **Phase 7A**: AWS::Serverless::Connector transformer with permission profiles for all service pairs
- [x] **Phase 7B-7C**: AWS::Serverless::Application (nested stacks) and AWS::Serverless::GraphQLApi (AppSync) transformers
- [x] **Phase 8**: Plugin system with Globals, ImplicitApi, ImplicitHttpApi, PolicyTemplates, DefaultDefinitionBody plugins
- [x] **Phase 9A**: Main translator orchestration with resource ordering and plugin lifecycle
- [x] **Phase 9B**: Command-line interface with full CLI options
- [x] **Phase 10**: Comprehensive test suite with 2,583 fixtures, benchmarks, and Python comparison tool

### Implementation Complete

All phases of the implementation are now complete. The translator is feature-complete and passes all test fixtures.

## Installation

```bash
go get github.com/lex00/aws-sam-translator-go
```

## CLI Usage

The `sam-translate` CLI tool transforms SAM templates to CloudFormation:

```bash
# Build the CLI
go build -o sam-translate ./cmd/sam-translate

# Transform a SAM template and write to file
sam-translate -t template.yaml -o output.yaml

# Transform and print to stdout
sam-translate -t template.yaml --stdout

# With verbose output
sam-translate -t template.yaml -o output.yaml --verbose

# Specify AWS region for partition detection
sam-translate -t template.yaml -o output.yaml --region us-gov-west-1

# Show help
sam-translate --help

# Show version
sam-translate --version
```

### CLI Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--template-file` | `-t` | Path to SAM template file (required) |
| `--output-template` | `-o` | Path to output CloudFormation template |
| `--stdout` | | Write output to stdout |
| `--verbose` | | Enable verbose logging |
| `--region` | | AWS region for partition detection |
| `--help` | `-h` | Show help message |
| `--version` | | Show version information |

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Transform error (invalid template, file not found) |
| 2 | Invalid arguments |

## Library Usage

### Template Transformation

```go
import "github.com/lex00/aws-sam-translator-go/pkg/translator"

// Create a new translator with default options
tr := translator.New()

// Or with custom options
tr := translator.NewWithOptions(translator.Options{
    Region:    "us-west-2",
    AccountID: "123456789012",
    StackName: "my-sam-app",
    Partition: "aws",
})

// Transform raw YAML/JSON bytes
input, err := os.ReadFile("template.yaml")
if err != nil {
    log.Fatal(err)
}

output, err := tr.TransformBytes(input)
if err != nil {
    log.Fatal(err)
}

os.WriteFile("output.json", output, 0644)
```

### Intrinsic Function Resolution

```go
import (
    "github.com/lex00/aws-sam-translator-go/pkg/intrinsics"
    "github.com/lex00/aws-sam-translator-go/pkg/types"
)

// Create a resolve context with default pseudo-parameters
ctx := intrinsics.NewResolveContext(template)

// Or with custom options
ctx := intrinsics.NewResolveContextWithOptions(template, intrinsics.ResolveContextOptions{
    AccountId: "123456789012",
    Region:    "us-west-2",
    StackName: "my-stack",
})

// Set parameter values
ctx.SetParameter("Environment", "production")

// Create registry and resolve intrinsics
registry := intrinsics.NewRegistry()
result, err := registry.Resolve(ctx, map[string]interface{}{
    "Fn::Sub": "arn:aws:s3:::${BucketName}-${AWS::Region}",
})

// Check for AWS::NoValue
if intrinsics.IsNoValue(result) {
    // Property should be removed
}
```

### Policy Template Expansion

```go
import "github.com/lex00/aws-sam-translator-go/pkg/policy"

// Create processor (loads 81 embedded templates)
p, err := policy.New()
if err != nil {
    log.Fatal(err)
}

// List available templates
for _, name := range p.TemplateNames() {
    fmt.Println(name)
}

// Expand a template with parameters
definition, err := p.Expand("DynamoDBCrudPolicy", map[string]interface{}{
    "TableName": "MyTable",
})

// Get just the IAM Statement array
statements, err := p.ExpandStatements("S3ReadPolicy", map[string]interface{}{
    "BucketName": "my-bucket",
})
```

## Goals

- Transform SAM templates to CloudFormation without Python runtime
- Single binary distribution
- Native integration with [cfn-lint-go](https://github.com/lex00/cfn-lint-go)

## Documentation

- [Architecture Overview](docs/ARCHITECTURE.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Testing Guide](docs/TESTING.md)
- [Research & Feasibility Analysis](docs/RESEARCH.md)
- [Changelog](CHANGELOG.md)

## License

Apache-2.0
