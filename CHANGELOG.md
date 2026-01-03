# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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
