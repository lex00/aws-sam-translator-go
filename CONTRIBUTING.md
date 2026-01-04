# Contributing to aws-sam-translator-go

Thank you for your interest in contributing to aws-sam-translator-go! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.23 or later
- golangci-lint for linting

### Getting Started

```bash
# Clone the repository
git clone https://github.com/lex00/aws-sam-translator-go.git
cd aws-sam-translator-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run ./...
```

## Project Structure

```
aws-sam-translator-go/
├── cmd/sam-translate/     # CLI entry point
├── pkg/
│   ├── cloudformation/    # CloudFormation resource models
│   ├── errors/            # Custom error types
│   ├── intrinsics/        # Intrinsic function handlers
│   ├── model/             # SAM and CF type definitions
│   │   ├── eventsources/  # Event source handlers
│   │   │   ├── pull/      # Pull-based (SQS, Kinesis, etc.)
│   │   │   └── push/      # Push-based (S3, SNS, API, etc.)
│   │   ├── iam/           # IAM models
│   │   └── lambda/        # Lambda models
│   ├── parser/            # YAML/JSON template parsing
│   ├── plugins/           # Plugin system
│   ├── policy/            # Policy template processor
│   ├── region/            # Region/partition configuration
│   ├── sam/               # SAM resource transformers
│   ├── translator/        # ID and ARN generation
│   ├── types/             # Core type definitions
│   └── utils/             # Utility functions
├── testdata/              # Test fixtures (2,583 files)
│   ├── input/             # SAM template inputs
│   └── output/            # Expected CF outputs
└── docs/                  # Documentation
```

## Adding a New SAM Resource Transformer

1. **Create the transformer file** in `pkg/sam/`:
   ```go
   // pkg/sam/myresource.go
   package sam

   type MyResourceTransformer struct {
       // dependencies
   }

   func NewMyResourceTransformer() *MyResourceTransformer {
       return &MyResourceTransformer{}
   }

   func (t *MyResourceTransformer) Transform(resource *types.Resource) ([]types.Resource, error) {
       // Transform SAM resource to CloudFormation resources
   }
   ```

2. **Create tests** in `pkg/sam/myresource_test.go`

3. **Update CHANGELOG.md** with your changes

## Adding a New Event Source Handler

### Push Event Sources (S3, SNS, API Gateway, etc.)

Create in `pkg/model/eventsources/push/`:

```go
// pkg/model/eventsources/push/myevent.go
package push

type MyEvent struct {
    // Event properties from SAM spec
}

func (e *MyEvent) ToCloudFormation(function *lambda.Function) ([]interface{}, error) {
    // Generate CloudFormation resources (Lambda::Permission, etc.)
}
```

### Pull Event Sources (SQS, Kinesis, DynamoDB Streams, etc.)

Create in `pkg/model/eventsources/pull/`:

```go
// pkg/model/eventsources/pull/myevent.go
package pull

type MyEvent struct {
    // Event properties from SAM spec
}

func (e *MyEvent) ToEventSourceMapping(function *lambda.Function) (*lambda.EventSourceMapping, error) {
    // Generate Lambda::EventSourceMapping
}
```

## Adding a New Plugin

1. Create plugin in `pkg/plugins/`:
   ```go
   // pkg/plugins/myplugin.go
   package plugins

   type MyPlugin struct{}

   func (p *MyPlugin) Name() string {
       return "MyPlugin"
   }

   func (p *MyPlugin) Priority() int {
       return 100 // Lower = runs first
   }

   func (p *MyPlugin) BeforeTransform(template *types.Template) error {
       // Modify template before SAM transformation
       return nil
   }

   func (p *MyPlugin) AfterTransform(template *types.Template) error {
       // Modify template after SAM transformation
       return nil
   }
   ```

2. Register in the plugin registry

## Code Style

### General Guidelines

- Follow standard Go conventions and idioms
- Use `gofmt` and `goimports` for formatting
- All exported types and functions must have godoc comments
- Avoid global state; prefer dependency injection

### Linting

The project uses golangci-lint with strict settings:

```bash
golangci-lint run ./...
```

Key lint rules:
- `errcheck`: All errors must be handled, including type assertions
- `govet`: Shadow variable detection enabled
- `staticcheck`: Static analysis checks
- `misspell`: Spelling errors in comments

### Testing

- Write table-driven tests where appropriate
- Include both success and error cases
- Test partition-specific behavior (aws, aws-cn, aws-us-gov)
- Match test file naming: `myfile.go` → `myfile_test.go`

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "basic case",
            input: "foo",
            want:  "bar",
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Test Fixtures

The `testdata/` directory contains 2,583 test fixtures ported from the upstream Python aws-sam-translator:

- `testdata/input/` - SAM template inputs (775 files)
- `testdata/output/` - Expected CloudFormation outputs
  - Default partition (778 files)
  - `aws-cn` partition (515 files)
  - `aws-us-gov` partition (515 files)

See [docs/TESTING.md](docs/TESTING.md) for detailed information.

## Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Adding or updating tests
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `chore`: Build process or auxiliary tool changes

Examples:
```
feat(sam): add AWS::Serverless::Function transformer
fix(intrinsics): handle nested Fn::Sub correctly
docs: update CHANGELOG with Phase 5A completion
test(connector): add tests for Lambda→DynamoDB permissions
```

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes with appropriate tests
3. Ensure all tests pass: `go test ./...`
4. Ensure linting passes: `golangci-lint run ./...`
5. Update CHANGELOG.md with your changes
6. Submit a pull request with a clear description

## Reporting Issues

When reporting issues, please include:

1. Go version (`go version`)
2. Operating system
3. Steps to reproduce
4. Expected vs actual behavior
5. Relevant SAM template (if applicable)

## License

By contributing, you agree that your contributions will be licensed under the Apache-2.0 License.
