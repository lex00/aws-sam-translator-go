package sam

import (
	"testing"
)

func TestNewConnectorTransformer(t *testing.T) {
	transformer := NewConnectorTransformer()
	if transformer == nil {
		t.Fatal("expected non-nil transformer")
	}
	if transformer.profiles == nil {
		t.Fatal("expected non-nil profiles")
	}
}

func TestConnectorTransformer_Transform_LambdaToDynamoDB(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
				"Runtime": "nodejs18.x",
			},
		},
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
			"Properties": map[string]interface{}{
				"TableName": "my-table",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Destination: ConnectorEndpoint{
			ID: "MyTable",
		},
		Permissions: []string{"Read", "Write"},
	}

	resources, err := transformer.Transform("MyConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a managed policy
	policyKey := "MyConnectorPolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected policy resource with key %q, got keys: %v", policyKey, getKeys(resources))
	}

	if policy["Type"] != "AWS::IAM::ManagedPolicy" {
		t.Errorf("expected Type 'AWS::IAM::ManagedPolicy', got %v", policy["Type"])
	}

	// Check metadata
	metadata, ok := policy["Metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("expected Metadata in policy")
	}
	connectorMeta, ok := metadata["aws:sam:connectors"].(map[string]interface{})
	if !ok {
		t.Fatal("expected aws:sam:connectors metadata")
	}
	myConnectorMeta, ok := connectorMeta["MyConnector"].(map[string]interface{})
	if !ok {
		t.Fatal("expected MyConnector in metadata")
	}
	sourceMeta := myConnectorMeta["Source"].(map[string]interface{})
	if sourceMeta["Type"] != "AWS::Serverless::Function" {
		t.Errorf("expected source type AWS::Serverless::Function, got %v", sourceMeta["Type"])
	}

	// Check properties
	props := policy["Properties"].(map[string]interface{})
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})

	if len(statements) < 1 {
		t.Fatal("expected at least one statement")
	}

	// Check roles
	roles := props["Roles"].([]interface{})
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}

	// Role should be a Ref to the function's role
	roleRef, ok := roles[0].(map[string]interface{})
	if !ok {
		t.Fatal("expected role to be a Ref")
	}
	if roleRef["Ref"] != "MyFunctionRole" {
		t.Errorf("expected Ref 'MyFunctionRole', got %v", roleRef["Ref"])
	}
}

func TestConnectorTransformer_Transform_SNSToLambda(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyTopic": map[string]interface{}{
			"Type": "AWS::SNS::Topic",
			"Properties": map[string]interface{}{
				"TopicName": "my-topic",
			},
		},
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
				"Runtime": "nodejs18.x",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyTopic",
		},
		Destination: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("SNSConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a Lambda permission
	var permKey string
	for key := range resources {
		if key == "SNSConnectorWriteLambdaPermission" {
			permKey = key
			break
		}
	}
	if permKey == "" {
		t.Fatalf("expected Lambda permission resource, got keys: %v", getKeys(resources))
	}

	perm := resources[permKey].(map[string]interface{})
	if perm["Type"] != "AWS::Lambda::Permission" {
		t.Errorf("expected Type 'AWS::Lambda::Permission', got %v", perm["Type"])
	}

	props := perm["Properties"].(map[string]interface{})
	if props["Action"] != "lambda:InvokeFunction" {
		t.Errorf("expected Action 'lambda:InvokeFunction', got %v", props["Action"])
	}
	if props["Principal"] != "sns.amazonaws.com" {
		t.Errorf("expected Principal 'sns.amazonaws.com', got %v", props["Principal"])
	}

	// Check SourceArn is a Ref to the topic
	sourceArn := props["SourceArn"].(map[string]interface{})
	if sourceArn["Ref"] != "MyTopic" {
		t.Errorf("expected SourceArn Ref to 'MyTopic', got %v", sourceArn)
	}
}

func TestConnectorTransformer_Transform_EventsRuleToSQS(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyRule": map[string]interface{}{
			"Type": "AWS::Events::Rule",
			"Properties": map[string]interface{}{
				"ScheduleExpression": "rate(1 minute)",
			},
		},
		"MyQueue": map[string]interface{}{
			"Type": "AWS::SQS::Queue",
			"Properties": map[string]interface{}{
				"QueueName": "my-queue",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyRule",
		},
		Destination: ConnectorEndpoint{
			ID: "MyQueue",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("EventsConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a queue policy
	policyKey := "EventsConnectorQueuePolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected queue policy resource, got keys: %v", getKeys(resources))
	}

	if policy["Type"] != "AWS::SQS::QueuePolicy" {
		t.Errorf("expected Type 'AWS::SQS::QueuePolicy', got %v", policy["Type"])
	}

	props := policy["Properties"].(map[string]interface{})

	// Check Queues
	queues := props["Queues"].([]interface{})
	if len(queues) != 1 {
		t.Fatalf("expected 1 queue, got %d", len(queues))
	}
	queueRef := queues[0].(map[string]interface{})
	if queueRef["Ref"] != "MyQueue" {
		t.Errorf("expected queue Ref to 'MyQueue', got %v", queueRef)
	}

	// Check policy document
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})
	if len(statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(statements))
	}

	stmt := statements[0].(map[string]interface{})
	if stmt["Action"] != "sqs:SendMessage" {
		t.Errorf("expected Action 'sqs:SendMessage', got %v", stmt["Action"])
	}
}

func TestConnectorTransformer_Transform_EventsRuleToSNS(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyRule": map[string]interface{}{
			"Type": "AWS::Events::Rule",
			"Properties": map[string]interface{}{
				"ScheduleExpression": "rate(1 minute)",
			},
		},
		"MyTopic": map[string]interface{}{
			"Type": "AWS::SNS::Topic",
			"Properties": map[string]interface{}{
				"TopicName": "my-topic",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyRule",
		},
		Destination: ConnectorEndpoint{
			ID: "MyTopic",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("EventsSNSConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a topic policy
	policyKey := "EventsSNSConnectorTopicPolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected topic policy resource, got keys: %v", getKeys(resources))
	}

	if policy["Type"] != "AWS::SNS::TopicPolicy" {
		t.Errorf("expected Type 'AWS::SNS::TopicPolicy', got %v", policy["Type"])
	}

	props := policy["Properties"].(map[string]interface{})

	// Check Topics
	topics := props["Topics"].([]interface{})
	if len(topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(topics))
	}

	// Check policy document
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})
	if len(statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(statements))
	}

	stmt := statements[0].(map[string]interface{})
	if stmt["Action"] != "sns:Publish" {
		t.Errorf("expected Action 'sns:Publish', got %v", stmt["Action"])
	}
}

func TestConnectorTransformer_Transform_WithExplicitTypes(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{}

	// Use explicit types without relying on template lookup
	connector := &Connector{
		Source: ConnectorEndpoint{
			Type: "AWS::Lambda::Function",
			Arn:  "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
			RoleName: map[string]interface{}{
				"Ref": "MyFunctionRole",
			},
		},
		Destination: ConnectorEndpoint{
			Type: "AWS::DynamoDB::Table",
			Arn:  "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable",
		},
		Permissions: []string{"Read"},
	}

	resources, err := transformer.Transform("ExplicitConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}
}

func TestConnectorTransformer_Transform_DuplicatePermissions(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
			"Properties": map[string]interface{}{
				"TableName": "my-table",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Destination: ConnectorEndpoint{
			ID: "MyTable",
		},
		Permissions: []string{"Read", "Read", "Write", "Write"},
	}

	resources, err := transformer.Transform("DupeConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should deduplicate and create single policy
	if len(resources) != 1 {
		t.Errorf("expected 1 resource (consolidated policy), got %d", len(resources))
	}
}

func TestConnectorTransformer_Transform_Error_MissingSource(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "NonExistent",
		},
		Destination: ConnectorEndpoint{
			ID: "MyTable",
		},
		Permissions: []string{"Read"},
	}

	_, err := transformer.Transform("BadConnector", connector, templateResources)
	if err == nil {
		t.Error("expected error for missing source resource")
	}
}

func TestConnectorTransformer_Transform_Error_MissingDestination(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Destination: ConnectorEndpoint{
			ID: "NonExistent",
		},
		Permissions: []string{"Read"},
	}

	_, err := transformer.Transform("BadConnector", connector, templateResources)
	if err == nil {
		t.Error("expected error for missing destination resource")
	}
}

func TestConnectorTransformer_Transform_Error_UnsupportedProfile(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyQueue": map[string]interface{}{
			"Type": "AWS::SQS::Queue",
		},
		"MyBucket": map[string]interface{}{
			"Type": "AWS::S3::Bucket",
		},
	}

	// SQS -> S3 is not a supported combination
	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyQueue",
		},
		Destination: ConnectorEndpoint{
			ID: "MyBucket",
		},
		Permissions: []string{"Write"},
	}

	_, err := transformer.Transform("UnsupportedConnector", connector, templateResources)
	if err == nil {
		t.Error("expected error for unsupported profile")
	}
}

func TestConnectorTransformer_TransformEmbedded(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
			"Properties": map[string]interface{}{
				"TableName": "my-table",
			},
		},
	}

	connectors := map[string]EmbeddedConnector{
		"TableConnector": {
			Properties: EmbeddedConnectorProperties{
				Destination: ConnectorEndpoint{
					ID: "MyTable",
				},
				Permissions: []string{"Read", "Write"},
			},
		},
	}

	resources, err := transformer.TransformEmbedded("MyFunction", "AWS::Serverless::Function", connectors, templateResources)
	if err != nil {
		t.Fatalf("TransformEmbedded failed: %v", err)
	}

	// Should create a policy with logical ID: MyFunctionTableConnectorPolicy
	policyKey := "MyFunctionTableConnectorPolicy"
	_, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Errorf("expected policy resource with key %q, got keys: %v", policyKey, getKeys(resources))
	}
}

func TestConnectorTransformer_TransformEmbedded_MultipleConnectors(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
		},
		"MyQueue": map[string]interface{}{
			"Type": "AWS::SQS::Queue",
		},
	}

	connectors := map[string]EmbeddedConnector{
		"TableConn": {
			Properties: EmbeddedConnectorProperties{
				Destination: ConnectorEndpoint{
					ID: "MyTable",
				},
				Permissions: []string{"Read"},
			},
		},
		"QueueConn": {
			Properties: EmbeddedConnectorProperties{
				Destination: ConnectorEndpoint{
					ID: "MyQueue",
				},
				Permissions: []string{"Write"},
			},
		},
	}

	resources, err := transformer.TransformEmbedded("MyFunction", "AWS::Serverless::Function", connectors, templateResources)
	if err != nil {
		t.Fatalf("TransformEmbedded failed: %v", err)
	}

	// Should create policies for both connectors
	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d: %v", len(resources), getKeys(resources))
	}
}

func TestExtractEmbeddedConnectors(t *testing.T) {
	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
			"Connectors": map[string]interface{}{
				"TableConnector": map[string]interface{}{
					"Properties": map[string]interface{}{
						"Destination": map[string]interface{}{
							"Id": "MyTable",
						},
						"Permissions": []interface{}{"Read", "Write"},
					},
				},
			},
		},
		"MyTable": map[string]interface{}{
			"Type": "AWS::DynamoDB::Table",
		},
	}

	result := ExtractEmbeddedConnectors(templateResources)

	funcConnectors, ok := result["MyFunction"]
	if !ok {
		t.Fatal("expected connectors for MyFunction")
	}

	tableConn, ok := funcConnectors["TableConnector"]
	if !ok {
		t.Fatal("expected TableConnector")
	}

	if tableConn.Properties.Destination.ID != "MyTable" {
		t.Errorf("expected destination ID 'MyTable', got %v", tableConn.Properties.Destination.ID)
	}

	if len(tableConn.Properties.Permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(tableConn.Properties.Permissions))
	}
}

func TestExtractEmbeddedConnectors_NoConnectors(t *testing.T) {
	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
	}

	result := ExtractEmbeddedConnectors(templateResources)

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestExtractEmbeddedConnectors_WithDestinationType(t *testing.T) {
	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Connectors": map[string]interface{}{
				"ExternalConnector": map[string]interface{}{
					"Properties": map[string]interface{}{
						"Destination": map[string]interface{}{
							"Type": "AWS::DynamoDB::Table",
							"Arn":  "arn:aws:dynamodb:us-east-1:123456789012:table/ExternalTable",
						},
						"Permissions": []interface{}{"Read"},
					},
				},
			},
		},
	}

	result := ExtractEmbeddedConnectors(templateResources)

	funcConnectors := result["MyFunction"]
	extConn := funcConnectors["ExternalConnector"]

	if extConn.Properties.Destination.Type != "AWS::DynamoDB::Table" {
		t.Errorf("expected destination type 'AWS::DynamoDB::Table', got %v", extConn.Properties.Destination.Type)
	}
}

func TestConnectorProfiles_GetProfile(t *testing.T) {
	profiles := NewConnectorProfiles()

	tests := []struct {
		name         string
		sourceType   string
		destType     string
		expectNil    bool
		expectedType string
	}{
		{
			name:         "Lambda to DynamoDB",
			sourceType:   TypeLambdaFunction,
			destType:     TypeDynamoDBTable,
			expectNil:    false,
			expectedType: "AWS::IAM::ManagedPolicy",
		},
		{
			name:         "Serverless Function to DynamoDB",
			sourceType:   TypeServerlessFunction,
			destType:     TypeDynamoDBTable,
			expectNil:    false,
			expectedType: "AWS::IAM::ManagedPolicy",
		},
		{
			name:         "SNS to Lambda",
			sourceType:   TypeSNSTopic,
			destType:     TypeLambdaFunction,
			expectNil:    false,
			expectedType: "AWS::Lambda::Permission",
		},
		{
			name:         "SNS to Serverless Function",
			sourceType:   TypeSNSTopic,
			destType:     TypeServerlessFunction,
			expectNil:    false,
			expectedType: "AWS::Lambda::Permission",
		},
		{
			name:         "Events Rule to SQS",
			sourceType:   TypeEventsRule,
			destType:     TypeSQSQueue,
			expectNil:    false,
			expectedType: "AWS::SQS::QueuePolicy",
		},
		{
			name:         "Events Rule to SNS",
			sourceType:   TypeEventsRule,
			destType:     TypeSNSTopic,
			expectNil:    false,
			expectedType: "AWS::SNS::TopicPolicy",
		},
		{
			name:       "Unsupported combination",
			sourceType: TypeSQSQueue,
			destType:   TypeS3Bucket,
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := profiles.GetProfile(tt.sourceType, tt.destType)
			if tt.expectNil {
				if profile != nil {
					t.Errorf("expected nil profile, got %v", profile)
				}
				return
			}

			if profile == nil {
				t.Fatal("expected non-nil profile")
			}

			if profile.ResourceType != tt.expectedType {
				t.Errorf("expected ResourceType %q, got %q", tt.expectedType, profile.ResourceType)
			}
		})
	}
}

func TestConnectorProfile_GetActions(t *testing.T) {
	profiles := NewConnectorProfiles()

	// Lambda -> DynamoDB
	profile := profiles.GetProfile(TypeLambdaFunction, TypeDynamoDBTable)
	if profile == nil {
		t.Fatal("expected profile for Lambda -> DynamoDB")
	}

	readActions := profile.GetActions("Read", TypeLambdaFunction, TypeDynamoDBTable)
	if len(readActions) == 0 {
		t.Error("expected Read actions for Lambda -> DynamoDB")
	}

	writeActions := profile.GetActions("Write", TypeLambdaFunction, TypeDynamoDBTable)
	if len(writeActions) == 0 {
		t.Error("expected Write actions for Lambda -> DynamoDB")
	}

	// Unknown permission type should return nil
	unknownActions := profile.GetActions("Unknown", TypeLambdaFunction, TypeDynamoDBTable)
	if unknownActions != nil {
		t.Errorf("expected nil for unknown permission, got %v", unknownActions)
	}
}

func TestConnectorProfile_GetPrincipal(t *testing.T) {
	profiles := NewConnectorProfiles()

	tests := []struct {
		name              string
		sourceType        string
		destType          string
		expectedPrincipal string
	}{
		{
			name:              "SNS to Lambda",
			sourceType:        TypeSNSTopic,
			destType:          TypeLambdaFunction,
			expectedPrincipal: "sns.amazonaws.com",
		},
		{
			name:              "S3 to Lambda",
			sourceType:        TypeS3Bucket,
			destType:          TypeLambdaFunction,
			expectedPrincipal: "s3.amazonaws.com",
		},
		{
			name:              "Events Rule to Lambda",
			sourceType:        TypeEventsRule,
			destType:          TypeLambdaFunction,
			expectedPrincipal: "events.amazonaws.com",
		},
		{
			name:              "API Gateway to Lambda",
			sourceType:        TypeAPIGatewayRestApi,
			destType:          TypeLambdaFunction,
			expectedPrincipal: "apigateway.amazonaws.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := profiles.GetProfile(tt.sourceType, tt.destType)
			if profile == nil {
				t.Fatalf("expected profile for %s -> %s", tt.sourceType, tt.destType)
			}

			principal := profile.GetPrincipal(tt.sourceType)
			if principal != tt.expectedPrincipal {
				t.Errorf("expected principal %q, got %q", tt.expectedPrincipal, principal)
			}
		})
	}
}

func TestConnectorProfile_GetResources(t *testing.T) {
	profiles := NewConnectorProfiles()

	destArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyTable", "Arn"},
	}

	// Lambda -> DynamoDB should have patterns for table and indexes
	profile := profiles.GetProfile(TypeLambdaFunction, TypeDynamoDBTable)
	if profile == nil {
		t.Fatal("expected profile")
	}

	resources := profile.GetResources("Read", destArn, nil, TypeDynamoDBTable, TypeLambdaFunction)
	if len(resources) != 2 {
		t.Errorf("expected 2 resources (table + indexes), got %d", len(resources))
	}

	// First should be direct ARN
	firstRes, ok := resources[0].(map[string]interface{})
	if !ok {
		t.Errorf("expected first resource to be a map")
	}
	getAtt, ok := firstRes["Fn::GetAtt"].([]interface{})
	if !ok || len(getAtt) < 2 || getAtt[0] != "MyTable" {
		t.Errorf("expected first resource to reference MyTable")
	}

	// Second should be Fn::Sub pattern for indexes
	fnSub, ok := resources[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected second resource to be Fn::Sub")
	}
	if _, hasSub := fnSub["Fn::Sub"]; !hasSub {
		t.Error("expected Fn::Sub in second resource")
	}
}

func TestNormalizeResourceType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{TypeServerlessFunction, TypeLambdaFunction},
		{TypeServerlessStateMachine, TypeStepFunctionsStateMachine},
		{TypeServerlessApi, TypeAPIGatewayRestApi},
		{TypeServerlessHttpApi, TypeAPIGatewayV2Api},
		{TypeLambdaFunction, TypeLambdaFunction},
		{TypeDynamoDBTable, TypeDynamoDBTable},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeResourceType(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeResourceType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConnectorTransformer_StepFunctionsToLambda(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyStateMachine": map[string]interface{}{
			"Type": "AWS::Serverless::StateMachine",
			"Properties": map[string]interface{}{
				"DefinitionUri": "statemachine.asl.json",
			},
		},
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Lambda::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyStateMachine",
		},
		Destination: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("SFNConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a managed policy for Step Functions to invoke Lambda
	policyKey := "SFNConnectorPolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected policy resource, got keys: %v", getKeys(resources))
	}

	if policy["Type"] != "AWS::IAM::ManagedPolicy" {
		t.Errorf("expected Type 'AWS::IAM::ManagedPolicy', got %v", policy["Type"])
	}

	props := policy["Properties"].(map[string]interface{})
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})

	if len(statements) == 0 {
		t.Fatal("expected at least one statement")
	}

	stmt := statements[0].(map[string]interface{})
	actions := stmt["Action"].([]interface{})

	// Should have lambda:InvokeFunction and lambda:InvokeAsync
	hasInvoke := false
	for _, a := range actions {
		if a == "lambda:InvokeFunction" || a == "lambda:InvokeAsync" {
			hasInvoke = true
			break
		}
	}
	if !hasInvoke {
		t.Error("expected lambda:InvokeFunction or lambda:InvokeAsync action")
	}
}

func TestConnectorTransformer_LambdaToS3(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
		"MyBucket": map[string]interface{}{
			"Type": "AWS::S3::Bucket",
			"Properties": map[string]interface{}{
				"BucketName": "my-bucket",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Destination: ConnectorEndpoint{
			ID: "MyBucket",
		},
		Permissions: []string{"Read", "Write"},
	}

	resources, err := transformer.Transform("S3Connector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	policyKey := "S3ConnectorPolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected policy resource, got keys: %v", getKeys(resources))
	}

	props := policy["Properties"].(map[string]interface{})
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})

	// Should have both read and write statements
	if len(statements) == 0 {
		t.Fatal("expected statements")
	}

	// Check that S3 actions are present
	hasS3Actions := false
	for _, s := range statements {
		stmt := s.(map[string]interface{})
		if actions, ok := stmt["Action"].([]interface{}); ok {
			for _, a := range actions {
				if aStr, ok := a.(string); ok {
					if len(aStr) > 3 && aStr[:3] == "s3:" {
						hasS3Actions = true
						break
					}
				}
			}
		}
	}
	if !hasS3Actions {
		t.Error("expected S3 actions in policy")
	}
}

func TestConnectorTransformer_LambdaToSQS(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
		"MyQueue": map[string]interface{}{
			"Type": "AWS::SQS::Queue",
			"Properties": map[string]interface{}{
				"QueueName": "my-queue",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Destination: ConnectorEndpoint{
			ID: "MyQueue",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("SQSConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	policyKey := "SQSConnectorPolicy"
	policy, ok := resources[policyKey].(map[string]interface{})
	if !ok {
		t.Fatalf("expected policy resource, got keys: %v", getKeys(resources))
	}

	props := policy["Properties"].(map[string]interface{})
	policyDoc := props["PolicyDocument"].(map[string]interface{})
	statements := policyDoc["Statement"].([]interface{})

	if len(statements) == 0 {
		t.Fatal("expected statements")
	}

	// Check for SQS actions
	stmt := statements[0].(map[string]interface{})
	actions := stmt["Action"].([]interface{})

	hasSendMessage := false
	for _, a := range actions {
		if a == "sqs:SendMessage" {
			hasSendMessage = true
			break
		}
	}
	if !hasSendMessage {
		t.Error("expected sqs:SendMessage action")
	}
}

func TestConnectorTransformer_APIGatewayToLambda(t *testing.T) {
	transformer := NewConnectorTransformer()

	templateResources := map[string]interface{}{
		"MyApi": map[string]interface{}{
			"Type": "AWS::Serverless::Api",
			"Properties": map[string]interface{}{
				"StageName": "prod",
			},
		},
		"MyFunction": map[string]interface{}{
			"Type": "AWS::Serverless::Function",
			"Properties": map[string]interface{}{
				"Handler": "index.handler",
			},
		},
	}

	connector := &Connector{
		Source: ConnectorEndpoint{
			ID: "MyApi",
		},
		Destination: ConnectorEndpoint{
			ID: "MyFunction",
		},
		Permissions: []string{"Write"},
	}

	resources, err := transformer.Transform("APIConnector", connector, templateResources)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create a Lambda permission
	var permKey string
	for key := range resources {
		if key == "APIConnectorWriteLambdaPermission" {
			permKey = key
			break
		}
	}
	if permKey == "" {
		t.Fatalf("expected Lambda permission, got keys: %v", getKeys(resources))
	}

	perm := resources[permKey].(map[string]interface{})
	props := perm["Properties"].(map[string]interface{})

	if props["Principal"] != "apigateway.amazonaws.com" {
		t.Errorf("expected Principal 'apigateway.amazonaws.com', got %v", props["Principal"])
	}
}

// Helper function to get map keys for error messages
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
