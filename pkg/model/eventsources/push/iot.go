// Package push provides push event source handlers for AWS SAM.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/iot"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// IoTRuleEvent represents a SAM IoTRule event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-iotrule.html
//
// SAM Template Syntax:
//
//	Events:
//	  MyIoTRule:
//	    Type: IoTRule
//	    Properties:
//	      Sql: "SELECT * FROM 'my/topic'"
//	      AwsIotSqlVersion: "2016-03-23"
//
// CloudFormation Resources Generated:
//   - AWS::IoT::TopicRule - The IoT rule with the SQL statement and Lambda action
//   - AWS::Lambda::Permission - Grants IoT permission to invoke the Lambda function
type IoTRuleEvent struct {
	// Sql is the SQL statement used to query the topic (required).
	// Example: "SELECT * FROM 'my/topic'"
	Sql interface{} `json:"Sql" yaml:"Sql"`

	// AwsIotSqlVersion is the version of the SQL rules engine to use (optional).
	// Valid values: "2015-10-08", "2016-03-23", "beta"
	// If not specified, the latest version is used.
	AwsIotSqlVersion interface{} `json:"AwsIotSqlVersion,omitempty" yaml:"AwsIotSqlVersion,omitempty"`
}

// IoTRuleEventSourceHandler handles IoTRule event sources.
type IoTRuleEventSourceHandler struct{}

// NewIoTRuleEventSourceHandler creates a new IoTRule event source handler.
func NewIoTRuleEventSourceHandler() *IoTRuleEventSourceHandler {
	return &IoTRuleEventSourceHandler{}
}

// GenerateResources generates CloudFormation resources for an IoTRule event source.
// It creates:
//  1. AWS::IoT::TopicRule - the IoT rule that triggers the Lambda function
//  2. AWS::Lambda::Permission - grants IoT permission to invoke the Lambda function
func (h *IoTRuleEventSourceHandler) GenerateResources(
	functionLogicalID string,
	eventLogicalID string,
	event *IoTRuleEvent,
) (map[string]interface{}, error) {
	if event.Sql == nil {
		return nil, fmt.Errorf("ioTRule event source requires a Sql property")
	}

	resources := make(map[string]interface{})

	// Generate logical IDs for the resources
	ruleLogicalID := fmt.Sprintf("%s%s", functionLogicalID, eventLogicalID)
	permissionLogicalID := fmt.Sprintf("%s%sPermission", functionLogicalID, eventLogicalID)

	// Build function ARN reference
	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{functionLogicalID, "Arn"},
	}

	// Create IoT TopicRule
	topicRule := iot.NewTopicRule(event.Sql, functionArn)

	// Add optional properties
	if event.AwsIotSqlVersion != nil {
		topicRule.WithAwsIotSqlVersion(event.AwsIotSqlVersion)
	}

	// Convert topic rule to CloudFormation format
	resources[ruleLogicalID] = topicRule.ToCloudFormation()

	// Create Lambda Permission
	// Build the rule ARN for the source ARN
	ruleArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{ruleLogicalID, "Arn"},
	}

	permission := lambda.NewIoTPermission(
		map[string]interface{}{"Ref": functionLogicalID},
		ruleArn,
	)

	// Convert permission to CloudFormation format
	resources[permissionLogicalID] = permission.ToCloudFormation()

	return resources, nil
}

// Validate validates the IoTRule event configuration.
func (h *IoTRuleEventSourceHandler) Validate(event *IoTRuleEvent) error {
	if event.Sql == nil {
		return fmt.Errorf("ioTRule event source requires a Sql property")
	}
	return nil
}
