# aws-sam-translator-go

A Go port of [aws-sam-translator](https://github.com/aws/serverless-application-model) for transforming AWS SAM templates to CloudFormation.

## Status

**Work in Progress** - See the [Implementation Roadmap](https://github.com/lex00/aws-sam-translator-go/issues/24) for current progress.

### Completed

- [x] **Phase 1C**: Intrinsic function handlers (Ref, Fn::Sub, Fn::GetAtt, Fn::FindInMap, pass-through handlers)
- [x] **Phase 2D**: Policy template processor with 81 SAM policy templates
- [x] **Phase 10** (partial): Test fixtures - 2,583 fixtures ported from upstream Python aws-sam-translator

### In Progress

- [ ] Phase 1A-1B: Core types, parser
- [ ] Phase 2A-2C: Intrinsics resolver, ID/ARN generators
- [ ] Phase 3-9: CloudFormation models, event sources, SAM transformers
- [ ] Phase 10: Remaining test suite (unit tests, Python comparison tool)

## Installation

```bash
go get github.com/lex00/aws-sam-translator-go
```

## Usage

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

- [Research & Feasibility Analysis](docs/RESEARCH.md)
- [Changelog](CHANGELOG.md)

## License

Apache-2.0
