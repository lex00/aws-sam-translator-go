// Package push provides event source handlers for push-based Lambda triggers.
// Push event sources actively invoke Lambda functions when events occur.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/s3"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// S3EventSource represents an S3 bucket event notification configuration for Lambda.
// This event source type is a "push" model where S3 actively invokes the Lambda function.
//
// Supported S3 events include:
//   - s3:ObjectCreated:*
//   - s3:ObjectCreated:Put
//   - s3:ObjectCreated:Post
//   - s3:ObjectCreated:Copy
//   - s3:ObjectCreated:CompleteMultipartUpload
//   - s3:ObjectRemoved:*
//   - s3:ObjectRemoved:Delete
//   - s3:ObjectRemoved:DeleteMarkerCreated
//   - s3:ObjectRestore:*
//   - s3:ObjectRestore:Post
//   - s3:ObjectRestore:Completed
//   - s3:ReducedRedundancyLostObject
//   - s3:Replication:*
//   - s3:LifecycleExpiration:*
//   - s3:LifecycleTransition
//   - s3:IntelligentTiering
//   - s3:ObjectTagging:*
//   - s3:ObjectAcl:Put
type S3EventSource struct {
	// Bucket is the S3 bucket name or ARN (required).
	// Can be a string or CloudFormation intrinsic function (Ref, GetAtt, etc.).
	Bucket interface{}

	// Events is the list of S3 bucket events to trigger on (required).
	// Examples: "s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"
	Events []string

	// Filter specifies optional filtering rules based on object key name.
	Filter *S3NotificationFilter
}

// S3NotificationFilter specifies filtering for S3 events based on object key patterns.
type S3NotificationFilter struct {
	// S3Key contains the filter rules for object key names.
	S3Key *S3KeyFilter
}

// S3KeyFilter contains filter rules for S3 object keys.
type S3KeyFilter struct {
	// Rules is a list of filter rules (prefix and/or suffix).
	Rules []S3FilterRule
}

// S3FilterRule represents a single filter rule for S3 object keys.
type S3FilterRule struct {
	// Name is the filter type. Valid values: "prefix" or "suffix".
	Name string

	// Value is the filter value to match.
	// Can be a string or CloudFormation intrinsic function.
	Value interface{}
}

// NewS3EventSource creates a new S3 event source with required parameters.
func NewS3EventSource(bucket interface{}, events []string) *S3EventSource {
	return &S3EventSource{
		Bucket: bucket,
		Events: events,
	}
}

// WithFilter adds filtering rules to the S3 event source.
func (s *S3EventSource) WithFilter(filter *S3NotificationFilter) *S3EventSource {
	s.Filter = filter
	return s
}

// WithPrefixFilter adds a prefix filter rule.
func (s *S3EventSource) WithPrefixFilter(prefix interface{}) *S3EventSource {
	if s.Filter == nil {
		s.Filter = &S3NotificationFilter{
			S3Key: &S3KeyFilter{
				Rules: []S3FilterRule{},
			},
		}
	}
	if s.Filter.S3Key == nil {
		s.Filter.S3Key = &S3KeyFilter{
			Rules: []S3FilterRule{},
		}
	}
	s.Filter.S3Key.Rules = append(s.Filter.S3Key.Rules, S3FilterRule{
		Name:  "prefix",
		Value: prefix,
	})
	return s
}

// WithSuffixFilter adds a suffix filter rule.
func (s *S3EventSource) WithSuffixFilter(suffix interface{}) *S3EventSource {
	if s.Filter == nil {
		s.Filter = &S3NotificationFilter{
			S3Key: &S3KeyFilter{
				Rules: []S3FilterRule{},
			},
		}
	}
	if s.Filter.S3Key == nil {
		s.Filter.S3Key = &S3KeyFilter{
			Rules: []S3FilterRule{},
		}
	}
	s.Filter.S3Key.Rules = append(s.Filter.S3Key.Rules, S3FilterRule{
		Name:  "suffix",
		Value: suffix,
	})
	return s
}

// ToCloudFormationResources generates CloudFormation resources for the S3 event source.
// This creates:
//   1. AWS::Lambda::Permission - Grants S3 permission to invoke the Lambda function
//   2. S3 bucket notification configuration (returned as metadata for bucket modification)
//
// The bucket property itself is not created, as the bucket is assumed to exist.
// The caller must apply the notification configuration to the bucket.
//
// Parameters:
//   - functionRef: Reference to the Lambda function (typically a Ref or GetAtt intrinsic)
//   - functionName: Logical name of the Lambda function for resource naming
//
// Returns:
//   - resources: Map of CloudFormation resources (Permission)
//   - notification: The LambdaConfiguration to add to the S3 bucket
//   - error: Any validation errors
func (s *S3EventSource) ToCloudFormationResources(functionRef interface{}, functionName string) (map[string]interface{}, *s3.LambdaConfiguration, error) {
	if s.Bucket == nil {
		return nil, nil, fmt.Errorf("S3 event source bucket is required")
	}
	if len(s.Events) == 0 {
		return nil, nil, fmt.Errorf("S3 event source must specify at least one event")
	}

	resources := make(map[string]interface{})

	// Create Lambda permission for S3 to invoke the function
	// The permission allows s3.amazonaws.com to invoke the function
	// SourceAccount is set to AWS::AccountId to prevent other accounts' buckets from invoking
	permission := lambda.NewS3Permission(
		functionRef,
		s.Bucket, // SourceArn - the bucket ARN
		map[string]interface{}{"Ref": "AWS::AccountId"}, // SourceAccount
	)

	// Generate a unique logical ID for the permission
	// Format: <FunctionName>S3Permission<hash>
	permissionID := fmt.Sprintf("%sS3Permission", functionName)

	resources[permissionID] = permission.ToCloudFormation()

	// Create S3 notification configurations for each event type
	// Multiple events can trigger the same function
	var notificationConfig *s3.LambdaConfiguration

	// For simplicity, we create one notification per event
	// In practice, SAM may consolidate multiple events into a single configuration
	// For now, we'll return the first event and the caller can handle multiple events
	if len(s.Events) > 0 {
		notificationConfig = &s3.LambdaConfiguration{
			Event:    s.Events[0], // Primary event
			Function: functionRef,
		}

		// Add filter if specified
		if s.Filter != nil && s.Filter.S3Key != nil && len(s.Filter.S3Key.Rules) > 0 {
			cfnFilter := &s3.NotificationFilter{
				S3Key: &s3.S3KeyFilter{
					Rules: make([]s3.FilterRule, len(s.Filter.S3Key.Rules)),
				},
			}
			for i, rule := range s.Filter.S3Key.Rules {
				cfnFilter.S3Key.Rules[i] = s3.FilterRule{
					Name:  rule.Name,
					Value: rule.Value,
				}
			}
			notificationConfig.Filter = cfnFilter
		}
	}

	return resources, notificationConfig, nil
}

// ToCloudFormationResourcesMultiEvent generates CloudFormation resources for S3 event source
// with support for multiple events.
//
// This variant creates separate notification configurations for each event type,
// which is more aligned with how SAM handles multiple S3 events.
//
// Parameters:
//   - functionRef: Reference to the Lambda function
//   - functionName: Logical name of the Lambda function
//
// Returns:
//   - resources: Map of CloudFormation resources (Permission)
//   - notifications: List of LambdaConfigurations for each event
//   - error: Any validation errors
func (s *S3EventSource) ToCloudFormationResourcesMultiEvent(functionRef interface{}, functionName string) (map[string]interface{}, []s3.LambdaConfiguration, error) {
	if s.Bucket == nil {
		return nil, nil, fmt.Errorf("S3 event source bucket is required")
	}
	if len(s.Events) == 0 {
		return nil, nil, fmt.Errorf("S3 event source must specify at least one event")
	}

	resources := make(map[string]interface{})

	// Create Lambda permission for S3 to invoke the function
	permission := lambda.NewS3Permission(
		functionRef,
		s.Bucket,
		map[string]interface{}{"Ref": "AWS::AccountId"},
	)

	permissionID := fmt.Sprintf("%sS3Permission", functionName)
	resources[permissionID] = permission.ToCloudFormation()

	// Create notification configuration for each event
	notifications := make([]s3.LambdaConfiguration, len(s.Events))
	for i, event := range s.Events {
		notifications[i] = s3.LambdaConfiguration{
			Event:    event,
			Function: functionRef,
		}

		// Add filter if specified (same filter applies to all events)
		if s.Filter != nil && s.Filter.S3Key != nil && len(s.Filter.S3Key.Rules) > 0 {
			cfnFilter := &s3.NotificationFilter{
				S3Key: &s3.S3KeyFilter{
					Rules: make([]s3.FilterRule, len(s.Filter.S3Key.Rules)),
				},
			}
			for j, rule := range s.Filter.S3Key.Rules {
				cfnFilter.S3Key.Rules[j] = s3.FilterRule{
					Name:  rule.Name,
					Value: rule.Value,
				}
			}
			notifications[i].Filter = cfnFilter
		}
	}

	return resources, notifications, nil
}

// Validate checks if the S3 event source configuration is valid.
func (s *S3EventSource) Validate() error {
	if s.Bucket == nil {
		return fmt.Errorf("S3 event source bucket is required")
	}
	if len(s.Events) == 0 {
		return fmt.Errorf("S3 event source must specify at least one event")
	}

	// Validate event names (basic validation)
	for _, event := range s.Events {
		if event == "" {
			return fmt.Errorf("S3 event name cannot be empty")
		}
	}

	// Validate filter rules if present
	if s.Filter != nil && s.Filter.S3Key != nil {
		for _, rule := range s.Filter.S3Key.Rules {
			if rule.Name != "prefix" && rule.Name != "suffix" {
				return fmt.Errorf("invalid S3 filter rule name: %s (must be 'prefix' or 'suffix')", rule.Name)
			}
			if rule.Value == nil {
				return fmt.Errorf("S3 filter rule value cannot be nil")
			}
		}
	}

	return nil
}

// GetBucket returns the bucket reference.
func (s *S3EventSource) GetBucket() interface{} {
	return s.Bucket
}

// GetEvents returns the list of events.
func (s *S3EventSource) GetEvents() []string {
	return s.Events
}

// GetFilter returns the notification filter.
func (s *S3EventSource) GetFilter() *S3NotificationFilter {
	return s.Filter
}
