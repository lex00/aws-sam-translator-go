package pull

import (
	"testing"
)

func TestNewCloudWatchLogsEventProperties(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	if props.LogGroupName != "/aws/lambda/my-function" {
		t.Errorf("expected LogGroupName '/aws/lambda/my-function', got %v", props.LogGroupName)
	}
	if props.FilterPattern != "ERROR" {
		t.Errorf("expected FilterPattern 'ERROR', got %s", props.FilterPattern)
	}
}

func TestCloudWatchLogsEventProperties_ToSubscriptionFilter_Minimal(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor")

	if sf.DestinationArn != "arn:aws:lambda:us-east-1:123456789012:function:my-processor" {
		t.Errorf("expected DestinationArn, got %v", sf.DestinationArn)
	}
	if sf.LogGroupName != "/aws/lambda/my-function" {
		t.Errorf("expected LogGroupName, got %v", sf.LogGroupName)
	}
	if sf.FilterPattern != "ERROR" {
		t.Errorf("expected FilterPattern 'ERROR', got %s", sf.FilterPattern)
	}
}

func TestCloudWatchLogsEventProperties_EmptyFilterPattern(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor")

	// Empty string is valid and matches all events
	if sf.FilterPattern != "" {
		t.Errorf("expected empty FilterPattern, got %s", sf.FilterPattern)
	}
}

func TestSubscriptionFilter_WithFilterName(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor").
		WithFilterName("MyErrorFilter")

	if sf.FilterName != "MyErrorFilter" {
		t.Errorf("expected FilterName 'MyErrorFilter', got %v", sf.FilterName)
	}
}

func TestSubscriptionFilter_WithRoleArn(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor").
		WithRoleArn("arn:aws:iam::123456789012:role/CWLtoLambdaRole")

	if sf.RoleArn != "arn:aws:iam::123456789012:role/CWLtoLambdaRole" {
		t.Errorf("expected RoleArn, got %v", sf.RoleArn)
	}
}

func TestSubscriptionFilter_WithDistribution(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor").
		WithDistribution("ByLogStream")

	if sf.Distribution != "ByLogStream" {
		t.Errorf("expected Distribution 'ByLogStream', got %s", sf.Distribution)
	}
}

func TestSubscriptionFilter_ToCloudFormation_Minimal(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"ERROR",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor")
	cfn := sf.ToCloudFormation()

	if cfn["Type"] != "AWS::Logs::SubscriptionFilter" {
		t.Errorf("expected Type 'AWS::Logs::SubscriptionFilter', got %v", cfn["Type"])
	}

	cfnProps := cfn["Properties"].(map[string]interface{})
	if cfnProps["DestinationArn"] != "arn:aws:lambda:us-east-1:123456789012:function:my-processor" {
		t.Errorf("expected DestinationArn in CFN properties")
	}
	if cfnProps["LogGroupName"] != "/aws/lambda/my-function" {
		t.Errorf("expected LogGroupName in CFN properties")
	}
	if cfnProps["FilterPattern"] != "ERROR" {
		t.Errorf("expected FilterPattern in CFN properties")
	}
}

func TestSubscriptionFilter_ToCloudFormation_FullConfig(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"[ERROR, WARN]",
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor").
		WithFilterName("MyFilter").
		WithRoleArn("arn:aws:iam::123456789012:role/CWLtoLambdaRole").
		WithDistribution("Random")

	cfn := sf.ToCloudFormation()
	cfnProps := cfn["Properties"].(map[string]interface{})

	if cfnProps["FilterName"] != "MyFilter" {
		t.Errorf("expected FilterName in CFN properties")
	}
	if cfnProps["RoleArn"] != "arn:aws:iam::123456789012:role/CWLtoLambdaRole" {
		t.Errorf("expected RoleArn in CFN properties")
	}
	if cfnProps["Distribution"] != "Random" {
		t.Errorf("expected Distribution in CFN properties")
	}
}

func TestSubscriptionFilter_ToCloudFormation_WithIntrinsicFunction(t *testing.T) {
	// Test with CloudFormation intrinsic function reference
	logGroupRef := map[string]interface{}{
		"Ref": "MyLogGroup",
	}
	destinationRef := map[string]interface{}{
		"Fn::GetAtt": []string{"MyFunction", "Arn"},
	}

	props := NewCloudWatchLogsEventProperties(logGroupRef, "")
	sf := props.ToSubscriptionFilter(destinationRef)
	cfn := sf.ToCloudFormation()

	cfnProps := cfn["Properties"].(map[string]interface{})

	// Verify intrinsic functions are preserved
	logGroupName := cfnProps["LogGroupName"].(map[string]interface{})
	if logGroupName["Ref"] != "MyLogGroup" {
		t.Errorf("expected LogGroupName Ref to be preserved")
	}

	destArn := cfnProps["DestinationArn"].(map[string]interface{})
	getAtt := destArn["Fn::GetAtt"].([]string)
	if getAtt[0] != "MyFunction" || getAtt[1] != "Arn" {
		t.Errorf("expected DestinationArn Fn::GetAtt to be preserved")
	}
}

func TestSubscriptionFilter_KinesisStreamDestination(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"[INFO, DEBUG]",
	)

	sf := props.ToSubscriptionFilter("arn:aws:kinesis:us-east-1:123456789012:stream/my-stream").
		WithRoleArn("arn:aws:iam::123456789012:role/CWLtoKinesisRole")

	cfn := sf.ToCloudFormation()
	cfnProps := cfn["Properties"].(map[string]interface{})

	if cfnProps["DestinationArn"] != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("expected Kinesis stream ARN as DestinationArn")
	}
	if cfnProps["RoleArn"] != "arn:aws:iam::123456789012:role/CWLtoKinesisRole" {
		t.Errorf("expected RoleArn for Kinesis destination")
	}
}

func TestSubscriptionFilter_FirehoseDestination(t *testing.T) {
	props := NewCloudWatchLogsEventProperties(
		"/aws/lambda/my-function",
		"",
	)

	sf := props.ToSubscriptionFilter("arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream").
		WithRoleArn("arn:aws:iam::123456789012:role/CWLtoFirehoseRole")

	cfn := sf.ToCloudFormation()
	cfnProps := cfn["Properties"].(map[string]interface{})

	if cfnProps["DestinationArn"] != "arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream" {
		t.Errorf("expected Firehose ARN as DestinationArn")
	}
}

func TestSubscriptionFilter_ComplexFilterPattern(t *testing.T) {
	// Test with a complex JSON filter pattern
	filterPattern := `{ $.eventType = "LOGIN" && $.result = "FAILURE" }`
	props := NewCloudWatchLogsEventProperties(
		"/aws/cloudtrail/my-trail",
		filterPattern,
	)

	sf := props.ToSubscriptionFilter("arn:aws:lambda:us-east-1:123456789012:function:my-processor")
	cfn := sf.ToCloudFormation()

	cfnProps := cfn["Properties"].(map[string]interface{})
	if cfnProps["FilterPattern"] != filterPattern {
		t.Errorf("expected complex filter pattern to be preserved")
	}
}
