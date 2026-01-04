# How To Use aws-sam-translator-go

This guide covers common usage patterns for transforming SAM templates to CloudFormation.

## Quick Start

### Installation

```bash
go install github.com/lex00/aws-sam-translator-go/cmd/sam-translate@latest
```

Or build from source:

```bash
git clone https://github.com/lex00/aws-sam-translator-go.git
cd aws-sam-translator-go
make build
```

### Basic Usage

Transform a SAM template to CloudFormation:

```bash
# Write output to a file
sam-translate --template-file template.yaml --output-template output.yaml

# Write output to stdout
sam-translate --template-file template.yaml --stdout

# Specify AWS region (affects ARN generation)
sam-translate --template-file template.yaml --region us-west-2 --stdout

# Verbose output for debugging
sam-translate --template-file template.yaml --verbose --stdout
```

## Library Usage

Use aws-sam-translator-go as a library in your Go application:

```go
package main

import (
    "fmt"
    "os"

    "github.com/lex00/aws-sam-translator-go/pkg/translator"
)

func main() {
    // Read SAM template
    input, err := os.ReadFile("template.yaml")
    if err != nil {
        panic(err)
    }

    // Create translator with options
    t := translator.NewWithOptions(translator.Options{
        Region:    "us-east-1",
        AccountID: "123456789012",
        StackName: "my-stack",
        Partition: "aws",
    })

    // Transform
    output, err := t.TransformBytes(input)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(output))
}
```

### Using the Parser Directly

```go
import "github.com/lex00/aws-sam-translator-go/pkg/parser"

p := parser.New()
template, err := p.Parse(yamlContent)
if err != nil {
    // Handle parse error
}

// Access template components
for name, resource := range template.Resources {
    fmt.Printf("Resource: %s, Type: %s\n", name, resource.Type)
}
```

### Using Policy Templates

```go
import "github.com/lex00/aws-sam-translator-go/pkg/policy"

// Create processor with embedded templates
processor, err := policy.New()
if err != nil {
    panic(err)
}

// List available templates
for _, name := range processor.TemplateNames() {
    fmt.Println(name)
}

// Expand a policy template
policyDoc, err := processor.Expand("DynamoDBCrudPolicy", map[string]interface{}{
    "TableName": "MyTable",
})

// Get just the statements array
statements, err := processor.ExpandStatements("S3ReadPolicy", map[string]interface{}{
    "BucketName": "my-bucket",
})
```

## SAM Template Examples

### Basic Function

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: ./src
      MemorySize: 128
      Timeout: 30
```

### Function with API Gateway

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Resources:
  HelloFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: ./src
      Events:
        HelloApi:
          Type: Api
          Properties:
            Path: /hello
            Method: GET
```

### Function with DynamoDB Event

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Resources:
  ProcessorFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: ./src
      Policies:
        - DynamoDBCrudPolicy:
            TableName: !Ref MyTable
      Events:
        Stream:
          Type: DynamoDB
          Properties:
            Stream: !GetAtt MyTable.StreamArn
            StartingPosition: TRIM_HORIZON
            BatchSize: 100

  MyTable:
    Type: AWS::Serverless::SimpleTable
    Properties:
      PrimaryKey:
        Name: id
        Type: String
```

### Using Globals

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Runtime: python3.9
    Timeout: 30
    MemorySize: 256
    Environment:
      Variables:
        LOG_LEVEL: INFO

Resources:
  Function1:
    Type: AWS::Serverless::Function
    Properties:
      Handler: function1.handler
      CodeUri: ./src

  Function2:
    Type: AWS::Serverless::Function
    Properties:
      Handler: function2.handler
      CodeUri: ./src
      Timeout: 60  # Overrides global
```

## Using Intrinsic Functions

### Ref

Reference parameters and resources:

```yaml
Parameters:
  Environment:
    Type: String
    Default: dev

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: ./src
      Environment:
        Variables:
          ENV: !Ref Environment
```

### Fn::Sub

String substitution with variables:

```yaml
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: !Sub 's3://${ArtifactBucket}/code.zip'
      Environment:
        Variables:
          TABLE_ARN: !Sub 'arn:aws:dynamodb:${AWS::Region}:${AWS::AccountId}:table/${TableName}'
```

### Fn::GetAtt

Get resource attributes:

```yaml
Resources:
  MyTable:
    Type: AWS::Serverless::SimpleTable

  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: ./src
      Environment:
        Variables:
          TABLE_NAME: !GetAtt MyTable.TableName
```

## Known Limitations

### ImportValue Restrictions

`Fn::ImportValue` has limited support in certain contexts:

- Cannot be used with `RestApiId` property
- Cannot be used with `Policies` property on Functions
- Cannot be used with `StageName` for APIs

### OpenAPI Generation

When using `AWS::Serverless::Api` or `AWS::Serverless::HttpApi` without a `DefinitionBody`, the translator automatically generates a complete OpenAPI specification from function events:

```yaml
Resources:
  HelloFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: ./src
      Events:
        HelloApi:
          Type: Api
          Properties:
            Path: /hello
            Method: GET
        CreateApi:
          Type: Api
          Properties:
            Path: /items
            Method: POST
```

The translator will automatically:
- Generate a Swagger 2.0 spec for `AWS::Serverless::Api`
- Generate an OpenAPI 3.0 spec for `AWS::Serverless::HttpApi`
- Add `x-amazon-apigateway-integration` extensions with Lambda proxy integration
- Extract path parameters from routes (e.g., `/users/{id}`)
- Merge routes into existing `DefinitionBody` if provided

For advanced configurations, you can still provide your own OpenAPI definition:

```yaml
Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      StageName: prod
      DefinitionBody:
        openapi: "3.0.1"
        info:
          title: "My API"
          version: "1.0"
        paths:
          /hello:
            get:
              x-amazon-apigateway-integration:
                type: aws_proxy
                httpMethod: POST
                uri: !Sub 'arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${HelloFunction.Arn}/invocations'
```

## Troubleshooting

### Template Parse Errors

If you get parse errors, check:
1. YAML syntax (proper indentation, no tabs)
2. Intrinsic function syntax (`!Ref` vs `Ref:`)
3. Required properties are present

### Transform Errors

Common transform errors:
- **Missing Handler**: Zip package type requires Handler property
- **Missing Runtime**: Required for Lambda functions
- **Invalid Event Type**: Check event type spelling and properties

Use verbose mode for detailed error information:

```bash
sam-translate --template-file template.yaml --verbose --stdout
```

## Integration with SAM CLI

This tool produces CloudFormation output compatible with `aws cloudformation deploy`:

```bash
# Transform SAM to CloudFormation
sam-translate --template-file template.yaml --output-template cfn-template.yaml

# Deploy with CloudFormation
aws cloudformation deploy \
  --template-file cfn-template.yaml \
  --stack-name my-stack \
  --capabilities CAPABILITY_IAM
```
