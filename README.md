# aws-sam-translator-go

A Go port of [aws-sam-translator](https://github.com/aws/serverless-application-model) for transforming AWS SAM templates to CloudFormation.

## Status

**Work in Progress** - See the [Implementation Roadmap](https://github.com/lex00/aws-sam-translator-go/issues/24) for current progress.

### Completed

- [x] **Phase 2D**: Policy template processor with 81 SAM policy templates

### In Progress

- [ ] Phase 1A-1C: Core types, parser, intrinsic handlers
- [ ] Phase 2A-2C: Intrinsics resolver, ID/ARN generators
- [ ] Phase 3-10: CloudFormation models, event sources, SAM transformers

## Installation

```bash
go get github.com/lex00/aws-sam-translator-go
```

## Usage

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
