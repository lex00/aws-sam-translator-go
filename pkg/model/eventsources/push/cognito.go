// Package push provides push event source handlers for AWS SAM.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// CognitoEvent represents a SAM Cognito event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-cognito.html
//
// SAM Template Syntax:
//
//	Events:
//	  CognitoUserPoolPreSignup:
//	    Type: Cognito
//	    Properties:
//	      UserPool: !Ref MyCognitoUserPool
//	      Trigger: PreSignUp
//
// CloudFormation Resources Generated:
//   - AWS::Lambda::Permission - Grants Cognito permission to invoke the Lambda function
//
// Note: The Cognito event source modifies the existing UserPool resource's LambdaConfig
// property rather than creating a new AWS::Cognito::UserPool resource.
type CognitoEvent struct {
	// UserPool is a reference to a Cognito UserPool defined in the same template (required).
	// This should be a Ref or GetAtt intrinsic function.
	UserPool interface{} `json:"UserPool" yaml:"UserPool"`

	// Trigger is the Cognito User Pool trigger type (required).
	// Valid values: PreSignUp, PostConfirmation, PreAuthentication, PostAuthentication,
	// PreTokenGeneration, CustomMessage, UserMigration, DefineAuthChallenge,
	// CreateAuthChallenge, VerifyAuthChallengeResponse
	Trigger interface{} `json:"Trigger" yaml:"Trigger"`
}

// CognitoTriggerType defines the valid Cognito trigger types.
type CognitoTriggerType string

const (
	CognitoTriggerPreSignUp                   CognitoTriggerType = "PreSignUp"
	CognitoTriggerPostConfirmation            CognitoTriggerType = "PostConfirmation"
	CognitoTriggerPreAuthentication           CognitoTriggerType = "PreAuthentication"
	CognitoTriggerPostAuthentication          CognitoTriggerType = "PostAuthentication"
	CognitoTriggerPreTokenGeneration          CognitoTriggerType = "PreTokenGeneration"
	CognitoTriggerCustomMessage               CognitoTriggerType = "CustomMessage"
	CognitoTriggerUserMigration               CognitoTriggerType = "UserMigration"
	CognitoTriggerDefineAuthChallenge         CognitoTriggerType = "DefineAuthChallenge"
	CognitoTriggerCreateAuthChallenge         CognitoTriggerType = "CreateAuthChallenge"
	CognitoTriggerVerifyAuthChallengeResponse CognitoTriggerType = "VerifyAuthChallengeResponse"
)

// ValidCognitoTriggers is a map of valid Cognito trigger types.
var ValidCognitoTriggers = map[CognitoTriggerType]bool{
	CognitoTriggerPreSignUp:                   true,
	CognitoTriggerPostConfirmation:            true,
	CognitoTriggerPreAuthentication:           true,
	CognitoTriggerPostAuthentication:          true,
	CognitoTriggerPreTokenGeneration:          true,
	CognitoTriggerCustomMessage:               true,
	CognitoTriggerUserMigration:               true,
	CognitoTriggerDefineAuthChallenge:         true,
	CognitoTriggerCreateAuthChallenge:         true,
	CognitoTriggerVerifyAuthChallengeResponse: true,
}

// CognitoEventSourceHandler handles Cognito event sources.
type CognitoEventSourceHandler struct{}

// NewCognitoEventSourceHandler creates a new Cognito event source handler.
func NewCognitoEventSourceHandler() *CognitoEventSourceHandler {
	return &CognitoEventSourceHandler{}
}

// GenerateResources generates CloudFormation resources for a Cognito event source.
// It creates:
//  1. AWS::Lambda::Permission - grants Cognito permission to invoke the Lambda function
//
// Additionally, it returns metadata for modifying the UserPool's LambdaConfig property.
func (h *CognitoEventSourceHandler) GenerateResources(
	functionLogicalID string,
	eventLogicalID string,
	event *CognitoEvent,
) (map[string]interface{}, error) {
	if event.UserPool == nil {
		return nil, fmt.Errorf("cognito event source requires a UserPool property")
	}
	if event.Trigger == nil {
		return nil, fmt.Errorf("cognito event source requires a Trigger property")
	}

	resources := make(map[string]interface{})

	// Generate logical ID for the permission
	permissionLogicalID := fmt.Sprintf("%s%sPermission", functionLogicalID, eventLogicalID)

	// Build the UserPool ARN for the source ARN
	// The UserPool property could be a Ref, GetAtt, or a direct ARN
	userPoolArn := h.buildUserPoolArn(event.UserPool)

	// Create Lambda Permission
	permission := lambda.NewCognitoPermission(
		map[string]interface{}{"Ref": functionLogicalID},
		userPoolArn,
	)

	// Convert permission to CloudFormation format
	resources[permissionLogicalID] = permission.ToCloudFormation()

	return resources, nil
}

// buildUserPoolArn constructs the UserPool ARN from the UserPool reference.
func (h *CognitoEventSourceHandler) buildUserPoolArn(userPool interface{}) interface{} {
	// If it's already an ARN string (starts with "arn:"), return as-is
	if arnStr, ok := userPool.(string); ok {
		if len(arnStr) > 4 && arnStr[:4] == "arn:" {
			return arnStr
		}
	}

	// If it's a Ref, convert to GetAtt for the ARN
	if refMap, ok := userPool.(map[string]interface{}); ok {
		if ref, hasRef := refMap["Ref"]; hasRef {
			return map[string]interface{}{
				"Fn::GetAtt": []interface{}{ref, "Arn"},
			}
		}
		// If it's already a GetAtt or other intrinsic, return as-is
		return refMap
	}

	// Default: wrap in GetAtt if it's a logical ID string
	return map[string]interface{}{
		"Fn::GetAtt": []interface{}{userPool, "Arn"},
	}
}

// GetLambdaConfigUpdate returns the LambdaConfig property update for the UserPool.
// This is used to modify the existing UserPool resource.
func (h *CognitoEventSourceHandler) GetLambdaConfigUpdate(
	functionLogicalID string,
	event *CognitoEvent,
) (string, interface{}, error) {
	trigger, ok := event.Trigger.(string)
	if !ok {
		return "", nil, fmt.Errorf("trigger must be a string, got %T", event.Trigger)
	}

	if !ValidCognitoTriggers[CognitoTriggerType(trigger)] {
		return "", nil, fmt.Errorf("invalid cognito trigger type: %s", trigger)
	}

	// Return the trigger name and the function ARN to set
	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{functionLogicalID, "Arn"},
	}

	return trigger, functionArn, nil
}

// Validate validates the Cognito event configuration.
func (h *CognitoEventSourceHandler) Validate(event *CognitoEvent) error {
	if event.UserPool == nil {
		return fmt.Errorf("cognito event source requires a UserPool property")
	}
	if event.Trigger == nil {
		return fmt.Errorf("cognito event source requires a Trigger property")
	}

	// Validate trigger type if it's a string
	if trigger, ok := event.Trigger.(string); ok {
		if !ValidCognitoTriggers[CognitoTriggerType(trigger)] {
			return fmt.Errorf("invalid cognito trigger type: %s", trigger)
		}
	}

	return nil
}
