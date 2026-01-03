# aws-sam-translator-go Research Document

## Executive Summary

Creating a Go port of aws-sam-translator is **highly feasible**. The Python codebase is well-structured with deterministic transformations. Estimated effort: **8-10 weeks for MVP, 12-16 weeks feature-complete**.

Key enabler: `cloudformation-schema-go` already exists and can provide CloudFormation resource type definitions.

---

## SAM Resource Types to Transform

| SAM Resource | CloudFormation Output | Complexity |
|--------------|----------------------|------------|
| AWS::Serverless::Function | Lambda, IAM Role, Permissions, Aliases | Very High |
| AWS::Serverless::Api | API Gateway REST, Swagger, Stages | Very High |
| AWS::Serverless::HttpApi | API Gateway HTTP, OpenAPI 3.0 | High |
| AWS::Serverless::StateMachine | Step Functions, IAM Role | High |
| AWS::Serverless::SimpleTable | DynamoDB Table | Low |
| AWS::Serverless::LayerVersion | Lambda Layer | Low |
| AWS::Serverless::Application | Nested Stack | Medium |
| AWS::Serverless::Connector | IAM Policies | High |
| AWS::Serverless::GraphQLApi | AppSync | High |

**Event Sources (17+ types):**
- Push: S3, SNS, API Gateway, IoT, Cognito, CloudWatch Events/Logs
- Pull: DynamoDB Streams, Kinesis, SQS, MSK, Kafka, MQ, DocumentDB
- Scheduled: EventBridge Schedule

---

## Architecture for Go Port

```
aws-sam-translator-go/
├── pkg/
│   ├── translator/        # Main transformation orchestrator
│   │   └── translator.go  # Process SAM → CloudFormation
│   ├── model/
│   │   ├── sam/           # SAM resource types
│   │   │   ├── function.go
│   │   │   ├── api.go
│   │   │   ├── httpapi.go
│   │   │   ├── statemachine.go
│   │   │   └── connector.go
│   │   ├── events/        # Event source handlers
│   │   │   ├── push.go    # S3, SNS, API Gateway, etc.
│   │   │   └── pull.go    # DynamoDB, Kinesis, SQS, etc.
│   │   └── cloudformation/ # CF resource wrappers
│   ├── plugins/           # Implicit API, policy templates
│   ├── openapi/           # Swagger/OpenAPI generation
│   └── policy/            # Policy template expansion
├── data/
│   └── policy_templates.json  # 30+ managed policies
├── cmd/
│   └── sam-translate/     # CLI tool
└── testdata/              # Test fixtures (port from Python)
```

---

## Key Dependencies

**Go Libraries:**
- `cloudformation-schema-go` - CloudFormation resource specs (exists!)
- `gopkg.in/yaml.v3` - YAML parsing with line tracking
- `github.com/go-playground/validator` - Struct validation
- `github.com/getkin/kin-openapi` - OpenAPI 3.0 handling

**Data Files to Port:**
- `policy_templates.json` - 30+ IAM policy templates
- Event source mappings
- Service permission matrices (for Connectors)

---

## Complexity Estimate

| Component | Estimated LOC |
|-----------|---------------|
| Parser/Validator | 500-800 |
| SAM Resource Models | 3,000-5,000 |
| Transformation Logic | 5,000-8,000 |
| Event Source Handlers | 1,500-2,000 |
| Connector Logic | 1,000-1,500 |
| OpenAPI Generation | 800-1,200 |
| Utilities | 2,000-3,000 |
| Tests | 5,000-8,000 |
| **Total** | **18,000-28,000** |

---

## Implementation Phases

**Phase 1: Foundation (2 weeks)**
- Project structure and build setup
- SAM template parser using yaml.v3
- Intrinsic function pass-through (!Ref, !GetAtt, !Sub)
- Integration with cloudformation-schema-go

**Phase 2: Core Transforms (2 weeks)**
- AWS::Serverless::Function → Lambda + IAM Role
- Basic event sources (API Gateway, S3, SNS)
- Policy template expansion

**Phase 3: API Support (2 weeks)**
- AWS::Serverless::Api → REST API Gateway
- AWS::Serverless::HttpApi → HTTP API Gateway
- OpenAPI/Swagger generation
- Implicit API plugin

**Phase 4: Advanced Resources (2 weeks)**
- AWS::Serverless::StateMachine
- AWS::Serverless::Connector
- Remaining event sources (Kinesis, SQS, DynamoDB Streams)

**Phase 5: Testing & Polish (2 weeks)**
- Port Python test fixtures (100+ templates)
- Output comparison with Python translator
- Performance benchmarking
- Documentation

---

## Critical Implementation Details

### 1. Intrinsic Function Handling

```go
// Must preserve, not evaluate
type IntrinsicRef struct {
    Ref string `yaml:"Ref" json:"Ref"`
}

type IntrinsicGetAtt struct {
    GetAtt []string `yaml:"Fn::GetAtt" json:"Fn::GetAtt"`
}

// Custom YAML unmarshaler to detect intrinsics
```

### 2. Resource Processing Order

```go
// Must process in this order:
// 1. Functions (API events modify API resources)
// 2. StateMachines (same reason)
// 3. APIs
// 4. Other resources
// 5. Connectors (need raw CF resources to exist)
```

### 3. Event Source Interface

```go
type EventSource interface {
    ToCloudFormation(function *SamFunction) ([]Resource, error)
    EventType() string
}

// Implementations: ApiEvent, S3Event, SQSEvent, etc.
```

### 4. Policy Template Expansion

```go
// policy_templates.json contains templates like:
// "SQSPollerPolicy": {
//   "Statement": [{
//     "Effect": "Allow",
//     "Action": ["sqs:DeleteMessage", "sqs:GetQueueAttributes", "sqs:ReceiveMessage"],
//     "Resource": {"Fn::Sub": "arn:${AWS::Partition}:sqs:${AWS::Region}:${AWS::AccountId}:${QueueName}"}
//   }]
// }
```

---

## Advantages of Go Implementation

1. **Single binary** - No Python runtime needed
2. **Fast startup** - Important for CI/CD pipelines
3. **Type safety** - Compile-time error checking
4. **Easy distribution** - Cross-platform binaries
5. **Integration** - Native with cfn-lint-go

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Output parity with Python | Comprehensive test suite comparing outputs |
| Missing edge cases | Port all 100+ Python test fixtures |
| OpenAPI complexity | Use existing Go OpenAPI libraries |
| Maintenance burden | Track aws-sam-translator releases |

---

## Validation Strategy

1. **Unit tests** - Each resource type transformation
2. **Integration tests** - Full template transformation
3. **Comparison tests** - Output matches Python version byte-for-byte
4. **Deployment tests** - Deploy transformed templates to AWS
5. **Performance tests** - Benchmark against Python version

---

## Success Criteria

- [ ] Transform all 10 SAM resource types
- [ ] Support all 17+ event sources
- [ ] 100% test fixture parity with Python
- [ ] <5% performance difference
- [ ] Deploy to production CloudFormation successfully
- [ ] Comprehensive documentation

---

## Recommended Repository Structure

**Option A: Standalone repo**
- `github.com/lex00/aws-sam-translator-go`
- Can be used by cfn-lint-go and other tools
- Independent release cycle

**Option B: Part of cfn-lint-go**
- `cfn-lint-go/pkg/samtranslator/`
- Tighter integration
- Single release cycle

**Recommendation: Option A** - Allows broader ecosystem adoption

---

## References

- [aws-sam-translator on PyPI](https://pypi.org/project/aws-sam-translator/) (~4.3M weekly downloads)
- [AWS SAM GitHub](https://github.com/aws/serverless-application-model)
- [SAM Specification](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html)
- [cloudformation-schema-go](https://github.com/lex00/cloudformation-schema-go)
