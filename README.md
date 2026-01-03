# aws-sam-translator-go

A Go port of [aws-sam-translator](https://github.com/aws/serverless-application-model) for transforming AWS SAM templates to CloudFormation.

## Status

**Work in Progress** - This project is in the research and planning phase.

See [docs/RESEARCH.md](docs/RESEARCH.md) for the feasibility analysis and implementation plan.

## Goals

- Transform SAM templates to CloudFormation without Python runtime
- Single binary distribution
- Native integration with [cfn-lint-go](https://github.com/lex00/cfn-lint-go)

## License

Apache-2.0
