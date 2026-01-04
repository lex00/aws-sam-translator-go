# Policy Templates Reference

Policy templates provide pre-configured IAM policy statements for common AWS service interactions. They offer scoped access instead of broad AWS managed policies.

## Why Use Policy Templates?

AWS managed policies (like `AmazonDynamoDBFullAccess`) grant access to all resources of that type. Policy templates provide least-privilege access to specific resources:

```yaml
# Too broad - access to ALL DynamoDB tables
Policies:
  - AmazonDynamoDBFullAccess

# Better - access only to MyTable
Policies:
  - DynamoDBCrudPolicy:
      TableName: !Ref MyTable
```

## Usage

Specify policy templates in the `Policies` property of `AWS::Serverless::Function`:

```yaml
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.9
      CodeUri: ./src
      Policies:
        - DynamoDBCrudPolicy:
            TableName: !Ref MyTable
        - S3ReadPolicy:
            BucketName: my-bucket
        - SQSPollerPolicy:
            QueueName: !GetAtt MyQueue.QueueName
```

Policy templates can be mixed with other policy types:

```yaml
Policies:
  # Policy template
  - DynamoDBCrudPolicy:
      TableName: !Ref MyTable
  # AWS managed policy ARN
  - arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess
  # Inline policy document
  - Version: '2012-10-17'
    Statement:
      - Effect: Allow
        Action: logs:*
        Resource: '*'
```

## Templates Without Parameters

Some templates require no parameters. Use an empty object:

```yaml
Policies:
  - CloudWatchPutMetricPolicy: {}
  - VPCAccessPolicy: {}
```

## Available Templates (81)

### DynamoDB

| Template | Parameters | Description |
|----------|------------|-------------|
| DynamoDBCrudPolicy | TableName | Full CRUD access to a DynamoDB table |
| DynamoDBReadPolicy | TableName | Read-only access to a DynamoDB table |
| DynamoDBWritePolicy | TableName | Write-only access to a DynamoDB table |
| DynamoDBStreamReadPolicy | TableName, StreamName | Read access to a DynamoDB stream |
| DynamoDBBackupFullAccessPolicy | TableName | Read/write access to DynamoDB on-demand backups |
| DynamoDBReconfigurePolicy | TableName | Permission to reconfigure a DynamoDB table |
| DynamoDBRestoreFromBackupPolicy | TableName | Permission to restore a table from backup |

### S3

| Template | Parameters | Description |
|----------|------------|-------------|
| S3CrudPolicy | BucketName | Full CRUD access to objects in an S3 bucket |
| S3ReadPolicy | BucketName | Read-only access to objects in an S3 bucket |
| S3WritePolicy | BucketName | Write-only access to objects in an S3 bucket |
| S3FullAccessPolicy | BucketName | Full access including bucket-level operations |

### SQS

| Template | Parameters | Description |
|----------|------------|-------------|
| SQSPollerPolicy | QueueName | Permission to poll an SQS queue |
| SQSSendMessagePolicy | QueueName | Permission to send messages to an SQS queue |

### SNS

| Template | Parameters | Description |
|----------|------------|-------------|
| SNSCrudPolicy | TopicName | Create, publish, and subscribe to SNS topics |
| SNSPublishMessagePolicy | TopicName | Permission to publish to an SNS topic |

### Kinesis

| Template | Parameters | Description |
|----------|------------|-------------|
| KinesisCrudPolicy | StreamName | Create, publish, and delete Kinesis streams |
| KinesisStreamReadPolicy | StreamName | List and read from a Kinesis stream |

### Kinesis Firehose

| Template | Parameters | Description |
|----------|------------|-------------|
| FirehoseCrudPolicy | DeliveryStreamName | Full access to a Firehose delivery stream |
| FirehoseWritePolicy | DeliveryStreamName | Write access to a Firehose delivery stream |

### Lambda

| Template | Parameters | Description |
|----------|------------|-------------|
| LambdaInvokePolicy | FunctionName | Permission to invoke a Lambda function |

### Step Functions

| Template | Parameters | Description |
|----------|------------|-------------|
| StepFunctionsExecutionPolicy | StateMachineName | Start state machine executions |
| StepFunctionsExecutionPolicy_v2 | StateMachineName | Start state machine executions (v2) |
| StepFunctionsCallbackPolicy | StateMachineName | Implement callback tasks in Step Functions |

### Secrets Manager

| Template | Parameters | Description |
|----------|------------|-------------|
| AWSSecretsManagerGetSecretValuePolicy | SecretArn | Get secret values from Secrets Manager |
| AWSSecretsManagerRotationPolicy | FunctionName | Rotate secrets in Secrets Manager |

### SSM Parameter Store

| Template | Parameters | Description |
|----------|------------|-------------|
| SSMParameterReadPolicy | ParameterName | Read a parameter from SSM Parameter Store |
| SSMParameterWithSlashPrefixReadPolicy | ParameterName | Read a parameter with slash prefix |

### KMS

| Template | Parameters | Description |
|----------|------------|-------------|
| KMSDecryptPolicy | KeyId | Decrypt with a KMS key |
| KMSEncryptPolicy | KeyId | Encrypt with a KMS key |
| KMSEncryptPolicy_v2 | KeyId | Encrypt with a KMS key (v2) |

### EventBridge

| Template | Parameters | Description |
|----------|------------|-------------|
| EventBridgePutEventsPolicy | EventBusName | Send events to EventBridge |

### CloudWatch

| Template | Parameters | Description |
|----------|------------|-------------|
| CloudWatchPutMetricPolicy | None | Put metrics to CloudWatch |
| CloudWatchDashboardPolicy | None | Operate on CloudWatch dashboards |
| CloudWatchDescribeAlarmHistoryPolicy | None | Describe CloudWatch alarm history |
| FilterLogEventsPolicy | LogGroupName | Filter log events from a log group |

### CodeCommit

| Template | Parameters | Description |
|----------|------------|-------------|
| CodeCommitCrudPolicy | RepositoryName | Full access to a CodeCommit repository |
| CodeCommitReadPolicy | RepositoryName | Read-only access to a CodeCommit repository |

### CodePipeline

| Template | Parameters | Description |
|----------|------------|-------------|
| CodePipelineLambdaExecutionPolicy | None | Lambda invoked by CodePipeline |
| CodePipelineReadOnlyPolicy | PipelineName | Read-only access to a pipeline |

### EC2

| Template | Parameters | Description |
|----------|------------|-------------|
| EC2DescribePolicy | None | Describe EC2 instances |
| EC2CopyImagePolicy | ImageId | Copy EC2 images |
| VPCAccessPolicy | None | Create/delete/describe ENIs for VPC |
| AMIDescribePolicy | None | Describe AMIs |

### ECS

| Template | Parameters | Description |
|----------|------------|-------------|
| EcsRunTaskPolicy | TaskDefinition | Start new ECS tasks |

### EKS

| Template | Parameters | Description |
|----------|------------|-------------|
| EKSDescribePolicy | None | Describe or list EKS clusters |

### EFS

| Template | Parameters | Description |
|----------|------------|-------------|
| EFSWriteAccessPolicy | FileSystem, AccessPoint | Mount EFS with write access |

### Elastic MapReduce (EMR)

| Template | Parameters | Description |
|----------|------------|-------------|
| ElasticMapReduceAddJobFlowStepsPolicy | ClusterId | Add steps to a running cluster |
| ElasticMapReduceCancelStepsPolicy | ClusterId | Cancel pending steps |
| ElasticMapReduceModifyInstanceFleetPolicy | ClusterId | Modify instance fleet capacities |
| ElasticMapReduceModifyInstanceGroupsPolicy | ClusterId | Modify instance group settings |
| ElasticMapReduceSetTerminationProtectionPolicy | ClusterId | Set termination protection |
| ElasticMapReduceTerminateJobFlowsPolicy | ClusterId | Terminate clusters |

### Elasticsearch/OpenSearch

| Template | Parameters | Description |
|----------|------------|-------------|
| ElasticsearchHttpPostPolicy | DomainName | POST and PUT to Elasticsearch |

### SES

| Template | Parameters | Description |
|----------|------------|-------------|
| SESCrudPolicy | IdentityName | Send email and verify identity |
| SESBulkTemplatedCrudPolicy | IdentityName | Send templated bulk emails |
| SESBulkTemplatedCrudPolicy_v2 | IdentityName, TemplateName | Send templated bulk emails (v2) |
| SESEmailTemplateCrudPolicy | None | Manage SES email templates |
| SESSendBouncePolicy | IdentityName | Send bounce notifications |

### Rekognition

| Template | Parameters | Description |
|----------|------------|-------------|
| RekognitionDetectOnlyPolicy | None | Detect faces, labels, and text |
| RekognitionFacesPolicy | None | Compare and detect faces |
| RekognitionLabelsPolicy | None | Detect object and moderation labels |
| RekognitionFacesManagementPolicy | CollectionId | Manage faces in a collection |
| RekognitionReadPolicy | CollectionId | List and search faces |
| RekognitionWriteOnlyAccessPolicy | CollectionId | Create collections and index faces |
| RekognitionNoDataAccessPolicy | CollectionId | Compare/detect without data access |

### Textract

| Template | Parameters | Description |
|----------|------------|-------------|
| TextractPolicy | None | Full Textract access |
| TextractDetectAnalyzePolicy | None | Detect and analyze documents |
| TextractGetResultPolicy | None | Get analyzed document results |

### Comprehend

| Template | Parameters | Description |
|----------|------------|-------------|
| ComprehendBasicAccessPolicy | None | Detect entities, sentiment, language |

### Polly

| Template | Parameters | Description |
|----------|------------|-------------|
| PollyFullAccessPolicy | LexiconName | Full access to Polly lexicon |

### SageMaker

| Template | Parameters | Description |
|----------|------------|-------------|
| SageMakerCreateEndpointPolicy | EndpointName | Create SageMaker endpoints |
| SageMakerCreateEndpointConfigPolicy | EndpointConfigName | Create endpoint configurations |

### Athena

| Template | Parameters | Description |
|----------|------------|-------------|
| AthenaQueryPolicy | WorkGroupName | Execute Athena queries |

### Pinpoint

| Template | Parameters | Description |
|----------|------------|-------------|
| PinpointEndpointAccessPolicy | PinpointApplicationId | Get/update Pinpoint endpoints |

### Route 53

| Template | Parameters | Description |
|----------|------------|-------------|
| Route53ChangeResourceRecordSetsPolicy | HostedZoneId | Change Route 53 record sets |

### ACM

| Template | Parameters | Description |
|----------|------------|-------------|
| AcmGetCertificatePolicy | CertificateArn | Retrieve ACM certificates |

### CloudFormation

| Template | Parameters | Description |
|----------|------------|-------------|
| CloudFormationDescribeStacksPolicy | None | Describe CloudFormation stacks |

### Organizations

| Template | Parameters | Description |
|----------|------------|-------------|
| OrganizationsListAccountsPolicy | None | List child accounts |

### Cost Explorer

| Template | Parameters | Description |
|----------|------------|-------------|
| CostExplorerReadOnlyPolicy | None | Read-only Cost Explorer access |

### Serverless Application Repository

| Template | Parameters | Description |
|----------|------------|-------------|
| ServerlessRepoReadWriteAccessPolicy | None | Create and list SAR applications |

### Mobile Analytics

| Template | Parameters | Description |
|----------|------------|-------------|
| MobileAnalyticsWriteOnlyAccessPolicy | None | Put event data for applications |

## Parameter Values

Parameters can use intrinsic functions:

```yaml
Policies:
  - DynamoDBCrudPolicy:
      TableName: !Ref MyTable  # Reference a resource
  - S3ReadPolicy:
      BucketName: !Sub '${AWS::StackName}-data'  # String substitution
  - LambdaInvokePolicy:
      FunctionName: !GetAtt OtherFunction.Arn  # Get attribute
```

## Expanded Policy Example

When you use `DynamoDBCrudPolicy`, the translator expands it to:

```json
{
  "PolicyDocument": {
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "dynamodb:GetItem",
          "dynamodb:DeleteItem",
          "dynamodb:PutItem",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:UpdateItem",
          "dynamodb:BatchWriteItem",
          "dynamodb:BatchGetItem",
          "dynamodb:DescribeTable",
          "dynamodb:ConditionCheckItem"
        ],
        "Resource": [
          "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable",
          "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable/index/*"
        ]
      }
    ]
  }
}
```
