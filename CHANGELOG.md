# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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
