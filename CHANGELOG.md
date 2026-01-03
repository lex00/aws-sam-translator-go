# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Logical ID and ARN Generators** (Phase 2B-2C - #5)
  - `pkg/translator/logical_id.go` - Hash-based LogicalIdGenerator with SHA256 for deterministic, stable CloudFormation logical IDs
  - `pkg/translator/arn.go` - Partition-aware ArnGenerator supporting aws, aws-cn, aws-us-gov partitions
  - `pkg/translator/verify_id.go` - IDVerifier for logical ID validation, duplicate detection, and ARN verification
  - Support for 15+ AWS services: Lambda, API Gateway, IAM, S3, DynamoDB, SNS, SQS, Kinesis, Step Functions, EventBridge, CloudWatch, Secrets Manager, KMS, Cognito, CodeDeploy
  - API deployment ID recalculation on OpenAPI spec change
  - IDStabilityChecker for deterministic ID generation verification

- **Intrinsics Resolver with Tree Traversal** (Phase 2A - #4)
  - `pkg/intrinsics/resolver.go` - Pre-order tree traversal resolver for nested intrinsic function resolution
  - `pkg/intrinsics/resource_refs.go` - DependencyTracker for building resource dependency graphs
  - Logical ID mutation support for SAM transformations
  - Placeholder protection for CloudFormation runtime intrinsics (Fn::ImportValue, Fn::GetAZs)
  - Topological sort for processing resources in dependency order
  - ResourceRefCollector for scanning templates for Ref and GetAtt references

- **YAML/JSON Parser with Intrinsic Support** (Phase 1B - #2)
  - `pkg/parser/yaml.go` - Custom YAML tag handling for all CloudFormation intrinsic short-form tags (!Ref, !Sub, !GetAtt, !Join, !If, !Select, !FindInMap, !Base64, !Cidr, !GetAZs, !ImportValue, !Split, !Transform, !And, !Equals, !Not, !Or, !Condition)
  - `pkg/parser/json.go` - JSON parsing with intrinsic structure preservation
  - `pkg/parser/intrinsics.go` - Intrinsic function detection, validation, and structure checking
  - `pkg/parser/parser.go` - Enhanced parser with source location tracking (line/column) for error reporting
  - Automatic conversion of !GetAtt dot notation (`!GetAtt Resource.Attr`) to array form
  - Template validation and auto-detection of YAML vs JSON format
  - 24 comprehensive tests

- **Core Types and Region Configuration** (Phase 1A - #1)
  - `pkg/region/config.go` - AWS partition detection for `aws`, `aws-cn`, and `aws-us-gov` partitions
  - Region-to-partition mapping functions (`GetPartitionForRegion`, `GetArnPartition`, `GetDNSSuffix`)
  - Region validation and default region handling
  - `pkg/utils/utils.go` - Extended with sorting functions (`SortedKeys`, `SortedResourceKeys`, `SortStringSlice`, `UniqueStrings`) for deterministic output
  - Comparison utilities (`DeepEqual`, `MapContains`, `StringSliceContains`, `MergeMaps`)
  - Safe accessor functions (`GetStringValue`, `GetMapValue`, `GetSliceValue`)
  - Comprehensive test coverage for errors, types, utils, and region packages

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
