# Development Guide

This guide helps contributors set up their development environment and understand the codebase conventions.

## Prerequisites

- Go 1.23 or later
- Make
- Git

## Environment Setup

```bash
# Clone the repository
git clone https://github.com/lex00/aws-sam-translator-go.git
cd aws-sam-translator-go

# Download dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for a specific package
go test ./pkg/translator/...

# Run a specific test
go test ./pkg/sam/... -run TestFunctionTransformer
```

## Development Guidelines

### 1. Do Not Resolve Runtime Intrinsics

Some CloudFormation intrinsic functions should be preserved for CloudFormation to resolve at deployment time. Only resolve intrinsics that can be evaluated at transform time:

**Resolve at transform time:**
- `Ref` (for parameters with default values)
- `Fn::Sub` (when all variables are known)
- `Fn::GetAtt` (for SAM-generated resources)
- `Fn::FindInMap` (when mapping exists)
- `Fn::Join` (when all values are strings)

**Preserve for CloudFormation:**
- `Fn::If`, `Fn::Select`, `Fn::Split`
- `Fn::ImportValue`, `Fn::GetAZs`
- `Fn::Base64`, `Fn::Cidr`
- Condition functions (`Fn::And`, `Fn::Or`, `Fn::Not`, `Fn::Equals`)

### 2. Maintain Output Parity with Python

The Go implementation should produce identical CloudFormation output as the Python aws-sam-translator for the same input. Use the test fixtures to verify:

```bash
# Compare output with Python version
scripts/compare-output.sh testdata/input/basic_function.yaml
```

### 3. Use Interface-Based Design

New components should implement interfaces for testability and extensibility:

```go
// Good: Interface-based design
type Transformer interface {
    Transform(logicalID string, resource Resource) (map[string]interface{}, error)
}

// Bad: Concrete type only
func TransformFunction(logicalID string, resource Resource) (map[string]interface{}, error)
```

### 4. Return Errors, Don't Panic

Always return errors instead of panicking:

```go
// Good
func ParseTemplate(data []byte) (*Template, error) {
    if len(data) == 0 {
        return nil, errors.New("empty template data")
    }
    // ...
}

// Bad
func ParseTemplate(data []byte) *Template {
    if len(data) == 0 {
        panic("empty template data")
    }
    // ...
}
```

## Code Conventions

### Package Organization

```
pkg/
├── types/          # Core data structures (Template, Resource, etc.)
├── parser/         # YAML/JSON parsing
├── intrinsics/     # Intrinsic function resolution
├── sam/            # SAM resource transformers
├── plugins/        # Plugin system
├── translator/     # Main orchestrator, ID/ARN generation
├── policy/         # Policy template processor
├── region/         # AWS region/partition utilities
├── model/          # CloudFormation resource models
│   ├── iam/
│   ├── lambda/
│   └── eventsources/
└── cloudformation/ # CloudFormation resource builders
```

### Naming Conventions

- Use `CamelCase` for exported functions and types
- Use `camelCase` for unexported functions and types
- Prefix interfaces with the noun they represent (e.g., `Transformer`, not `ITransformer`)
- Use descriptive names over abbreviations

### Error Handling

Use custom error types from `pkg/errors/`:

```go
import "github.com/lex00/aws-sam-translator-go/pkg/errors"

// For invalid resource configuration
return nil, errors.NewInvalidResourceException(
    logicalID,
    "Handler is required for Zip package type",
)

// For invalid event configuration
return nil, errors.NewInvalidEventException(
    logicalID,
    eventType,
    "Queue is required for SQS event",
)
```

### Documentation

Add godoc comments for all exported types and functions:

```go
// FunctionTransformer transforms AWS::Serverless::Function resources
// into CloudFormation Lambda resources with associated IAM roles,
// permissions, and event source mappings.
type FunctionTransformer struct {
    // ...
}

// Transform converts a SAM Function to CloudFormation resources.
// It returns a map of logical IDs to resources, or an error if
// the function configuration is invalid.
func (t *FunctionTransformer) Transform(logicalID string, fn *Function) (map[string]interface{}, error) {
    // ...
}
```

## Adding New Features

### Adding a New SAM Resource Type

1. Define the resource struct in `pkg/sam/`
2. Create a transformer implementing the transformation logic
3. Register the transformer in `pkg/translator/translator.go`
4. Add test fixtures in `testdata/input/` and expected output in `testdata/output/`
5. Update documentation

### Adding a New Event Source

1. Add the event type to `pkg/model/eventsources/push/` or `pull/`
2. Implement the `EventSource` interface
3. Register in the function transformer
4. Add test fixtures

### Adding a New Plugin

1. Implement the `Plugin` interface in `pkg/plugins/`
2. Register in the default plugins list
3. Add tests

## Profiling

Use Go's built-in profiling tools:

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./pkg/translator/...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./pkg/translator/...
go tool pprof mem.prof

# Benchmarks
go test -bench=. ./pkg/translator/...
```

## Verifying Transforms

To verify your changes produce correct output:

1. Run the full test suite: `make test`
2. Compare specific templates with Python output
3. Test with real SAM templates from AWS examples

```bash
# Transform a template
./bin/sam-translate --template-file template.yaml --stdout

# Compare with Python sam-cli
sam validate --template-file template.yaml
```

## Pull Request Checklist

Before submitting a PR:

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`make lint`)
- [ ] New code has tests
- [ ] Public APIs have godoc comments
- [ ] CHANGELOG.md is updated (for user-facing changes)
