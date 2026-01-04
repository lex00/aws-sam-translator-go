# Testing Guide

This document describes the testing structure and how to run tests for aws-sam-translator-go.

## Quick Start

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./pkg/sam/...

# Run a specific test
go test -v ./pkg/sam -run TestFunctionTransformer

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Organization

### Unit Tests

Each package has corresponding `*_test.go` files:

| Package | Test File | Description |
|---------|-----------|-------------|
| `pkg/sam` | `function_test.go` | Function transformer tests |
| `pkg/sam` | `statemachine_test.go` | StateMachine transformer tests |
| `pkg/sam` | `connector_test.go` | Connector transformer tests |
| `pkg/sam` | `simpletable_test.go` | SimpleTable transformer tests |
| `pkg/sam` | `layerversion_test.go` | LayerVersion transformer tests |
| `pkg/intrinsics` | `resolver_test.go` | Intrinsic resolver tests |
| `pkg/intrinsics` | `ref_test.go` | Ref action tests |
| `pkg/intrinsics` | `sub_test.go` | Fn::Sub action tests |
| `pkg/intrinsics` | `getatt_test.go` | Fn::GetAtt action tests |
| `pkg/intrinsics` | `findinmap_test.go` | Fn::FindInMap action tests |
| `pkg/parser` | `parser_test.go` | YAML/JSON parser tests |
| `pkg/policy` | `processor_test.go` | Policy template tests |
| `pkg/plugins` | `globals_test.go` | Globals plugin tests |
| `pkg/plugins` | `implicit_api_test.go` | Implicit REST API tests |
| `pkg/plugins` | `implicit_httpapi_test.go` | Implicit HTTP API tests |
| `pkg/translator` | `logical_id_test.go` | Logical ID generator tests |
| `pkg/translator` | `arn_test.go` | ARN generator tests |

### Event Source Tests

Push event sources (`pkg/model/eventsources/push/`):
- `s3_test.go` - S3 event tests
- `sns_test.go` - SNS event tests
- `api_test.go` - REST API Gateway tests
- `httpapi_test.go` - HTTP API Gateway tests
- `schedule_test.go` - Schedule event tests
- `cloudwatch_test.go` - CloudWatch Events tests
- `cognito_test.go` - Cognito trigger tests
- `iot_test.go` - IoT rule tests

Pull event sources (`pkg/model/eventsources/pull/`):
- `sqs_test.go` - SQS polling tests
- `kinesis_test.go` - Kinesis stream tests
- `dynamodb_test.go` - DynamoDB Streams tests
- `documentdb_test.go` - DocumentDB tests
- `msk_test.go` - MSK tests
- `mq_test.go` - Amazon MQ tests
- `cloudwatchlogs_test.go` - CloudWatch Logs tests
- `selfmanagedkafka_test.go` - Self-managed Kafka tests
- `schedulev2_test.go` - EventBridge Scheduler tests

## Test Fixtures

The `testdata/` directory contains **2,583 test fixtures** ported from the upstream Python [aws-sam-translator](https://github.com/aws/serverless-application-model).

### Directory Structure

```
testdata/
├── input/                    # 775 SAM template inputs
│   ├── success/             # 513 valid templates (expected to transform)
│   └── error/               # 262 invalid templates (expected to fail)
└── output/                  # Expected CloudFormation outputs
    ├── aws/                 # 778 outputs for default (aws) partition
    ├── aws-cn/              # 515 outputs for China partition
    └── aws-us-gov/          # 515 outputs for GovCloud partition
```

### Fixture Statistics

| Category | Count |
|----------|-------|
| Input SAM templates | 775 |
| Success cases | 513 |
| Error cases | 262 |
| Default partition outputs | 778 |
| China (aws-cn) outputs | 515 |
| GovCloud (aws-us-gov) outputs | 515 |
| **Total fixtures** | **2,583** |

### SAM Resource Coverage

The fixtures cover all SAM resource types:

- `AWS::Serverless::Function` - Lambda functions with all event types
- `AWS::Serverless::Api` - REST API Gateway
- `AWS::Serverless::HttpApi` - HTTP API Gateway
- `AWS::Serverless::SimpleTable` - DynamoDB tables
- `AWS::Serverless::LayerVersion` - Lambda layers
- `AWS::Serverless::StateMachine` - Step Functions
- `AWS::Serverless::Connector` - IAM policy connectors
- `AWS::Serverless::Application` - Nested stacks
- `AWS::Serverless::GraphQLApi` - AppSync APIs

### Event Source Coverage

Fixtures include all 18 event source types:

**Push Events:**
- S3 bucket notifications
- SNS topic subscriptions
- REST API Gateway (Api)
- HTTP API Gateway (HttpApi)
- Schedule (EventBridge)
- CloudWatch Events
- Cognito User Pool triggers
- IoT Rules

**Pull Events:**
- SQS queue polling
- Kinesis streams
- DynamoDB Streams
- DocumentDB change streams
- MSK (Managed Kafka)
- Amazon MQ
- CloudWatch Logs subscriptions
- Self-managed Kafka
- EventBridge Scheduler (v2)

### Using Fixtures in Tests

```go
func TestWithFixtures(t *testing.T) {
    // Read input SAM template
    input, err := os.ReadFile("testdata/input/success/function_basic.yaml")
    if err != nil {
        t.Fatal(err)
    }

    // Parse and transform
    template, err := parser.Parse(input)
    if err != nil {
        t.Fatal(err)
    }

    result, err := transformer.Transform(template)
    if err != nil {
        t.Fatal(err)
    }

    // Read expected output
    expected, err := os.ReadFile("testdata/output/aws/function_basic.json")
    if err != nil {
        t.Fatal(err)
    }

    // Compare
    if !reflect.DeepEqual(result, expected) {
        t.Errorf("output mismatch")
    }
}
```

### Partition-Specific Testing

Some resources generate different ARNs based on AWS partition. Test all three:

```go
partitions := []string{"aws", "aws-cn", "aws-us-gov"}

for _, partition := range partitions {
    t.Run(partition, func(t *testing.T) {
        ctx := intrinsics.NewResolveContextWithOptions(template,
            intrinsics.ResolveContextOptions{
                Partition: partition,
            })
        // Test partition-specific behavior
    })
}
```

## Test Patterns

### Table-Driven Tests

Most tests use the table-driven pattern:

```go
func TestLogicalIdGenerator(t *testing.T) {
    tests := []struct {
        name     string
        input    []string
        expected string
    }{
        {
            name:     "simple function",
            input:    []string{"MyFunction"},
            expected: "MyFunction",
        },
        {
            name:     "with suffix",
            input:    []string{"MyFunction", "Role"},
            expected: "MyFunctionRole",
        },
        {
            name:     "long name truncated",
            input:    []string{"VeryLongResourceNameThatExceedsLimit", "Role"},
            expected: "VeryLongResourceNa1a2b3c4Role", // truncated with hash
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gen := NewLogicalIdGenerator()
            got := gen.Generate(tt.input...)
            if got != tt.expected {
                t.Errorf("Generate() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Error Case Testing

Always test error conditions:

```go
func TestInvalidInput(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr string
    }{
        {
            name:    "missing required property",
            input:   `{"Type": "AWS::Serverless::Function"}`,
            wantErr: "Handler is required",
        },
        {
            name:    "invalid runtime",
            input:   `{"Type": "AWS::Serverless::Function", "Properties": {"Runtime": "invalid"}}`,
            wantErr: "invalid runtime",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := Transform(tt.input)
            if err == nil {
                t.Fatal("expected error, got nil")
            }
            if !strings.Contains(err.Error(), tt.wantErr) {
                t.Errorf("error = %v, want containing %v", err, tt.wantErr)
            }
        })
    }
}
```

## Linting

Run the linter before submitting:

```bash
golangci-lint run ./...
```

The `.golangci.yml` configuration enforces:

- `errcheck` - All errors handled (including type assertions)
- `govet` - Shadow variable detection
- `staticcheck` - Static analysis
- `ineffassign` - Ineffective assignments
- `unused` - Unused code detection
- `misspell` - Spelling errors
- `gofmt` / `goimports` - Formatting

## Continuous Integration

Tests run automatically on pull requests via GitHub Actions:

```yaml
# .github/workflows/test.yml
- name: Test
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Lint
  run: golangci-lint run ./...
```

## Debugging Tests

### Verbose Output

```bash
go test -v ./pkg/sam -run TestConnector
```

### Debug Specific Fixture

```bash
# Set environment variable to print debug info
DEBUG=1 go test -v ./pkg/sam -run TestFunctionWithS3Event
```

### Compare with Python Output

For parity testing with the Python aws-sam-translator:

```bash
# Generate output with Python translator
sam-translate --template testdata/input/function.yaml > python_output.json

# Generate output with Go translator
./sam-translate --template testdata/input/function.yaml > go_output.json

# Compare
diff python_output.json go_output.json
```
