# Globals Section Reference

The `Globals` section allows you to define common properties that all SAM resources of a given type inherit. This reduces repetition and ensures consistency across your serverless application.

## Basic Syntax

```yaml
Globals:
  Function:
    Runtime: python3.9
    Timeout: 30
    MemorySize: 256

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      CodeUri: ./src
      # Inherits Runtime, Timeout, MemorySize from Globals
```

## Supported Resource Types

The Globals section supports these SAM resource types:

- `Function` - AWS::Serverless::Function
- `Api` - AWS::Serverless::Api
- `HttpApi` - AWS::Serverless::HttpApi
- `SimpleTable` - AWS::Serverless::SimpleTable

## Function Globals

### Supported Properties

| Property | Type | Merge Behavior |
|----------|------|----------------|
| Handler | String | Override |
| Runtime | String | Override |
| CodeUri | String/S3Location | Override |
| DeadLetterQueue | Object | Override |
| Description | String | Override |
| MemorySize | Integer | Override |
| Timeout | Integer | Override |
| VpcConfig | Object | Override |
| Environment | Object | **Merge Variables** |
| Tags | Object | **Merge** |
| Tracing | String | Override |
| KmsKeyArn | String | Override |
| Layers | List | Override |
| AutoPublishAlias | String | Override |
| DeploymentPreference | Object | Override |
| PermissionsBoundary | String | Override |
| ReservedConcurrentExecutions | Integer | Override |
| ProvisionedConcurrencyConfig | Object | Override |
| AssumeRolePolicyDocument | Object | Override |
| EventInvokeConfig | Object | Override |

### Example

```yaml
Globals:
  Function:
    Runtime: python3.9
    Timeout: 30
    MemorySize: 256
    Tracing: Active
    Environment:
      Variables:
        LOG_LEVEL: INFO
        STAGE: production
    Tags:
      Team: backend
      Project: my-app

Resources:
  Function1:
    Type: AWS::Serverless::Function
    Properties:
      Handler: function1.handler
      CodeUri: ./src/function1
      # Inherits all globals

  Function2:
    Type: AWS::Serverless::Function
    Properties:
      Handler: function2.handler
      CodeUri: ./src/function2
      Timeout: 60  # Override global timeout
      Environment:
        Variables:
          FEATURE_FLAG: enabled  # Merged with global variables
```

## Api Globals

### Supported Properties

| Property | Type | Merge Behavior |
|----------|------|----------------|
| Auth | Object | Override |
| Name | String | Override |
| DefinitionUri | String | Override |
| CacheClusterEnabled | Boolean | Override |
| CacheClusterSize | String | Override |
| Variables | Object | Override |
| EndpointConfiguration | Object | Override |
| MethodSettings | List | Override |
| BinaryMediaTypes | List | Override |
| MinimumCompressionSize | Integer | Override |
| Cors | String/Object | Override |
| GatewayResponses | Object | Override |
| AccessLogSetting | Object | Override |
| CanarySetting | Object | Override |
| TracingEnabled | Boolean | Override |
| OpenApiVersion | String | Override |
| Domain | Object | Override |

### Example

```yaml
Globals:
  Api:
    Cors:
      AllowOrigin: "'*'"
      AllowHeaders: "'Content-Type,Authorization'"
      AllowMethods: "'GET,POST,PUT,DELETE'"
    Auth:
      DefaultAuthorizer: MyCognitoAuthorizer
      Authorizers:
        MyCognitoAuthorizer:
          UserPoolArn: !GetAtt MyUserPool.Arn
    TracingEnabled: true
    AccessLogSetting:
      DestinationArn: !GetAtt ApiLogGroup.Arn

Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      StageName: prod
      # Inherits Cors, Auth, TracingEnabled, AccessLogSetting
```

## HttpApi Globals

### Supported Properties

| Property | Type | Merge Behavior |
|----------|------|----------------|
| Auth | Object | Override |
| AccessLogSettings | Object | Override |
| StageVariables | Object | Override |
| Tags | Object | **Merge** |
| RouteSettings | Object | Override |
| FailOnWarnings | Boolean | Override |
| Domain | Object | Override |
| CorsConfiguration | Object | Override |

### Example

```yaml
Globals:
  HttpApi:
    Auth:
      DefaultAuthorizer: OAuth2Authorizer
      Authorizers:
        OAuth2Authorizer:
          AuthorizationScopes:
            - read
            - write
          JwtConfiguration:
            issuer: https://auth.example.com
            audience:
              - api
    CorsConfiguration:
      AllowOrigins:
        - https://example.com
      AllowMethods:
        - GET
        - POST
      AllowHeaders:
        - Authorization

Resources:
  MyHttpApi:
    Type: AWS::Serverless::HttpApi
    Properties:
      StageName: prod
      # Inherits Auth and CorsConfiguration
```

## SimpleTable Globals

### Supported Properties

| Property | Type | Merge Behavior |
|----------|------|----------------|
| SSESpecification | Object | Override |

### Example

```yaml
Globals:
  SimpleTable:
    SSESpecification:
      SSEEnabled: true

Resources:
  Table1:
    Type: AWS::Serverless::SimpleTable
    Properties:
      PrimaryKey:
        Name: id
        Type: String
      # Inherits SSESpecification

  Table2:
    Type: AWS::Serverless::SimpleTable
    Properties:
      PrimaryKey:
        Name: pk
        Type: String
      # Inherits SSESpecification
```

## Override Behavior

### Simple Override

Resource-level properties completely replace global properties:

```yaml
Globals:
  Function:
    Timeout: 30

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Timeout: 60  # Uses 60, not 30
```

### Merge Behavior

For `Environment.Variables` and `Tags`, values are merged:

```yaml
Globals:
  Function:
    Environment:
      Variables:
        GLOBAL_VAR: global_value
        SHARED_VAR: from_global
    Tags:
      Team: backend

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Environment:
        Variables:
          LOCAL_VAR: local_value
          SHARED_VAR: from_local  # Overrides global
      Tags:
        Service: my-service

# Result for MyFunction:
# Environment.Variables:
#   GLOBAL_VAR: global_value
#   SHARED_VAR: from_local
#   LOCAL_VAR: local_value
# Tags:
#   Team: backend
#   Service: my-service
```

## Properties NOT Supported in Globals

Some properties cannot be specified in Globals and must be set on each resource:

### Function
- `FunctionName`
- `Role`
- `Policies`
- `Events`

### Api/HttpApi
- `StageName`
- `DefinitionBody`

## Best Practices

### 1. Use Globals for Common Settings

```yaml
Globals:
  Function:
    Runtime: python3.9
    Timeout: 30
    Tracing: Active
```

### 2. Override Only When Necessary

```yaml
Resources:
  QuickFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: quick.handler
      # Uses global Timeout of 30

  SlowFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: slow.handler
      Timeout: 300  # Override for long-running function
```

### 3. Use Environment Variables for Shared Configuration

```yaml
Globals:
  Function:
    Environment:
      Variables:
        TABLE_NAME: !Ref DataTable
        QUEUE_URL: !Ref ProcessingQueue
```

### 4. Centralize Security Settings

```yaml
Globals:
  Function:
    VpcConfig:
      SecurityGroupIds:
        - !Ref LambdaSecurityGroup
      SubnetIds:
        - !Ref PrivateSubnet1
        - !Ref PrivateSubnet2
    KmsKeyArn: !GetAtt EncryptionKey.Arn
```

## Complete Example

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Runtime: python3.9
    Timeout: 30
    MemorySize: 256
    Tracing: Active
    Environment:
      Variables:
        LOG_LEVEL: INFO
        TABLE_NAME: !Ref DataTable
    Tags:
      Application: my-app
      Environment: production

  Api:
    Cors:
      AllowOrigin: "'*'"
    TracingEnabled: true

Resources:
  DataTable:
    Type: AWS::Serverless::SimpleTable

  ReadFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: read.handler
      CodeUri: ./src/read
      Policies:
        - DynamoDBReadPolicy:
            TableName: !Ref DataTable
      Events:
        Api:
          Type: Api
          Properties:
            Path: /data/{id}
            Method: GET

  WriteFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: write.handler
      CodeUri: ./src/write
      Timeout: 60  # Override for write operations
      Policies:
        - DynamoDBCrudPolicy:
            TableName: !Ref DataTable
      Events:
        Api:
          Type: Api
          Properties:
            Path: /data
            Method: POST
```
