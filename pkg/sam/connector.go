// Package sam provides SAM resource transformers.
package sam

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/model/iam"
)

// AWS resource type constants
const (
	TypeLambdaFunction            = "AWS::Lambda::Function"
	TypeServerlessFunction        = "AWS::Serverless::Function"
	TypeDynamoDBTable             = "AWS::DynamoDB::Table"
	TypeSNSTopic                  = "AWS::SNS::Topic"
	TypeSQSQueue                  = "AWS::SQS::Queue"
	TypeS3Bucket                  = "AWS::S3::Bucket"
	TypeStepFunctionsStateMachine = "AWS::StepFunctions::StateMachine"
	TypeServerlessStateMachine    = "AWS::Serverless::StateMachine"
	TypeEventsRule                = "AWS::Events::Rule"
	TypeAPIGatewayRestApi         = "AWS::ApiGateway::RestApi"
	TypeAPIGatewayV2Api           = "AWS::ApiGatewayV2::Api"
	TypeServerlessApi             = "AWS::Serverless::Api"
	TypeServerlessHttpApi         = "AWS::Serverless::HttpApi"
	TypeLocationPlaceIndex        = "AWS::Location::PlaceIndex"
	TypeAppSyncGraphQLApi         = "AWS::AppSync::GraphQLApi"
	TypeEventsEventBus            = "AWS::Events::EventBus"
)

// ConnectorEndpoint represents a source or destination in a connector.
type ConnectorEndpoint struct {
	// ID is the logical ID of the resource in the template.
	ID string `json:"Id,omitempty" yaml:"Id,omitempty"`

	// Type is the CloudFormation resource type.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// Arn is the ARN of the resource (can be intrinsic function).
	Arn interface{} `json:"Arn,omitempty" yaml:"Arn,omitempty"`

	// RoleName is the role name (for Lambda functions).
	RoleName interface{} `json:"RoleName,omitempty" yaml:"RoleName,omitempty"`

	// QueueUrl is the queue URL (for SQS).
	QueueUrl interface{} `json:"QueueUrl,omitempty" yaml:"QueueUrl,omitempty"`

	// Name is the resource name.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// ResourceId is used for API Gateway.
	ResourceId interface{} `json:"ResourceId,omitempty" yaml:"ResourceId,omitempty"`

	// Qualifier is used for Lambda aliases/versions.
	Qualifier interface{} `json:"Qualifier,omitempty" yaml:"Qualifier,omitempty"`
}

// Connector represents an AWS::Serverless::Connector resource.
type Connector struct {
	// Source is the source endpoint.
	Source ConnectorEndpoint `json:"Source" yaml:"Source"`

	// Destination is the destination endpoint.
	Destination ConnectorEndpoint `json:"Destination" yaml:"Destination"`

	// Permissions is the list of permissions (Read, Write).
	Permissions []string `json:"Permissions" yaml:"Permissions"`
}

// EmbeddedConnector represents a connector embedded in a resource's Connectors property.
type EmbeddedConnector struct {
	// Properties contains the connector properties.
	Properties EmbeddedConnectorProperties `json:"Properties" yaml:"Properties"`
}

// EmbeddedConnectorProperties contains the properties for an embedded connector.
type EmbeddedConnectorProperties struct {
	// Destination is the destination endpoint.
	Destination ConnectorEndpoint `json:"Destination" yaml:"Destination"`

	// Permissions is the list of permissions.
	Permissions []string `json:"Permissions" yaml:"Permissions"`
}

// ConnectorTransformer transforms AWS::Serverless::Connector resources.
type ConnectorTransformer struct {
	profiles *ConnectorProfiles
}

// NewConnectorTransformer creates a new ConnectorTransformer.
func NewConnectorTransformer() *ConnectorTransformer {
	return &ConnectorTransformer{
		profiles: NewConnectorProfiles(),
	}
}

// Transform converts a SAM Connector to CloudFormation resources.
// Returns a map of logical ID to CloudFormation resource.
func (t *ConnectorTransformer) Transform(logicalID string, connector *Connector, templateResources map[string]interface{}) (map[string]interface{}, error) {
	// Resolve source and destination types from template if ID is provided
	sourceType, err := t.resolveResourceType(connector.Source, templateResources)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve source type: %w", err)
	}

	destType, err := t.resolveResourceType(connector.Destination, templateResources)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve destination type: %w", err)
	}

	// Get the profile for this source/destination combination
	profile := t.profiles.GetProfile(sourceType, destType)
	if profile == nil {
		return nil, fmt.Errorf("no connector profile found for %s -> %s", sourceType, destType)
	}

	// Build the resources based on the profile
	resources := make(map[string]interface{})

	// Deduplicate permissions
	permissions := t.deduplicatePermissions(connector.Permissions)

	for _, perm := range permissions {
		switch profile.ResourceType {
		case "AWS::IAM::ManagedPolicy":
			policyResource, policyID := t.createManagedPolicy(logicalID, connector, sourceType, destType, perm, profile, templateResources)
			resources[policyID] = policyResource
		case "AWS::Lambda::Permission":
			permResource, permID := t.createLambdaPermission(logicalID, connector, sourceType, destType, perm, profile, templateResources)
			resources[permID] = permResource
		case "AWS::SQS::QueuePolicy":
			policyResource, policyID := t.createQueuePolicy(logicalID, connector, sourceType, destType, perm, profile, templateResources)
			resources[policyID] = policyResource
		case "AWS::SNS::TopicPolicy":
			policyResource, policyID := t.createTopicPolicy(logicalID, connector, sourceType, destType, perm, profile, templateResources)
			resources[policyID] = policyResource
		}
	}

	// If there are multiple policies with same resource type, consolidate them
	resources = t.consolidatePolicies(logicalID, resources, connector, sourceType, destType)

	return resources, nil
}

// TransformEmbedded transforms embedded connectors from a resource.
// sourceID is the logical ID of the parent resource.
// connectors is the map of connector names to EmbeddedConnector.
func (t *ConnectorTransformer) TransformEmbedded(sourceID, sourceType string, connectors map[string]EmbeddedConnector, templateResources map[string]interface{}) (map[string]interface{}, error) {
	allResources := make(map[string]interface{})

	for connectorName, embedded := range connectors {
		// Create a full Connector from the embedded connector
		connector := &Connector{
			Source: ConnectorEndpoint{
				ID:   sourceID,
				Type: sourceType,
			},
			Destination: embedded.Properties.Destination,
			Permissions: embedded.Properties.Permissions,
		}

		// Generate a logical ID for this embedded connector
		connectorLogicalID := sourceID + connectorName

		// Transform the connector
		resources, err := t.Transform(connectorLogicalID, connector, templateResources)
		if err != nil {
			return nil, fmt.Errorf("failed to transform embedded connector %s: %w", connectorName, err)
		}

		// Merge into allResources
		for id, resource := range resources {
			allResources[id] = resource
		}
	}

	return allResources, nil
}

// resolveResourceType gets the resource type from the endpoint or the template.
func (t *ConnectorTransformer) resolveResourceType(endpoint ConnectorEndpoint, templateResources map[string]interface{}) (string, error) {
	// If Type is explicitly provided, use it
	if endpoint.Type != "" {
		return endpoint.Type, nil
	}

	// If ID is provided, look up the resource in the template
	if endpoint.ID != "" {
		resource, ok := templateResources[endpoint.ID]
		if !ok {
			return "", fmt.Errorf("resource %q not found in template", endpoint.ID)
		}
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("resource %q is not a valid resource object", endpoint.ID)
		}
		resourceType, ok := resourceMap["Type"].(string)
		if !ok {
			return "", fmt.Errorf("resource %q does not have a Type", endpoint.ID)
		}
		return resourceType, nil
	}

	return "", fmt.Errorf("endpoint must have either Id or Type specified")
}

// deduplicatePermissions removes duplicate permissions.
func (t *ConnectorTransformer) deduplicatePermissions(permissions []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if !seen[p] {
			seen[p] = true
			result = append(result, p)
		}
	}
	return result
}

// createManagedPolicy creates an AWS::IAM::ManagedPolicy resource.
func (t *ConnectorTransformer) createManagedPolicy(
	logicalID string,
	connector *Connector,
	sourceType, destType, permission string,
	profile *ConnectorProfile,
	templateResources map[string]interface{},
) (map[string]interface{}, string) {
	// Get role reference for the source
	roleRef := t.getRoleReference(connector.Source, sourceType, templateResources)

	// Get destination ARN
	destArn := t.getResourceArn(connector.Destination, destType, templateResources)

	// Get source ARN for some profiles
	sourceArn := t.getResourceArn(connector.Source, sourceType, templateResources)

	// Build policy document
	policyDoc := iam.NewPolicyDocument()

	// Get actions for this permission type
	actions := profile.GetActions(permission, sourceType, destType)
	resources := profile.GetResources(permission, destArn, sourceArn, destType, sourceType)

	if len(actions) > 0 && len(resources) > 0 {
		stmt := iam.NewAllowStatement()
		stmt.Action = t.toInterfaceSlice(actions)
		stmt.Resource = resources
		policyDoc.AddStatement(stmt)
	}

	// Build metadata
	metadata := t.buildConnectorMetadata(logicalID, sourceType, destType)

	policyID := logicalID + "Policy"
	return map[string]interface{}{
		"Type":     "AWS::IAM::ManagedPolicy",
		"Metadata": metadata,
		"Properties": map[string]interface{}{
			"PolicyDocument": policyDoc.ToMap(),
			"Roles":          []interface{}{roleRef},
		},
	}, policyID
}

// createLambdaPermission creates an AWS::Lambda::Permission resource.
func (t *ConnectorTransformer) createLambdaPermission(
	logicalID string,
	connector *Connector,
	sourceType, destType, permission string,
	profile *ConnectorProfile,
	templateResources map[string]interface{},
) (map[string]interface{}, string) {
	// Get function ARN
	functionArn := t.getResourceArn(connector.Destination, destType, templateResources)

	// Get source ARN
	sourceArn := t.getSourceArnForLambdaPermission(connector.Source, sourceType, templateResources)

	// Get principal
	principal := profile.GetPrincipal(sourceType)

	// Build metadata
	metadata := t.buildConnectorMetadata(logicalID, sourceType, destType)

	permID := logicalID + permission + "LambdaPermission"
	resource := map[string]interface{}{
		"Type":     "AWS::Lambda::Permission",
		"Metadata": metadata,
		"Properties": map[string]interface{}{
			"Action":       "lambda:InvokeFunction",
			"FunctionName": functionArn,
			"Principal":    principal,
		},
	}

	// Add SourceArn if applicable
	if sourceArn != nil {
		if props, ok := resource["Properties"].(map[string]interface{}); ok {
			props["SourceArn"] = sourceArn
		}
	}

	return resource, permID
}

// createQueuePolicy creates an AWS::SQS::QueuePolicy resource.
func (t *ConnectorTransformer) createQueuePolicy(
	logicalID string,
	connector *Connector,
	sourceType, destType, permission string,
	profile *ConnectorProfile,
	templateResources map[string]interface{},
) (map[string]interface{}, string) {
	// Get queue URL
	queueUrl := t.getQueueUrl(connector.Destination, templateResources)
	queueArn := t.getResourceArn(connector.Destination, destType, templateResources)
	sourceArn := t.getResourceArn(connector.Source, sourceType, templateResources)

	// Build policy document
	policyDoc := iam.NewPolicyDocument()
	stmt := iam.NewAllowStatement()
	stmt.Action = "sqs:SendMessage"
	stmt.Resource = queueArn
	stmt.Principal = map[string]interface{}{"Service": profile.GetPrincipal(sourceType)}
	stmt.Condition = map[string]interface{}{
		"ArnEquals": map[string]interface{}{
			"aws:SourceArn": sourceArn,
		},
	}
	policyDoc.AddStatement(stmt)

	metadata := t.buildConnectorMetadata(logicalID, sourceType, destType)

	policyID := logicalID + "QueuePolicy"
	return map[string]interface{}{
		"Type":     "AWS::SQS::QueuePolicy",
		"Metadata": metadata,
		"Properties": map[string]interface{}{
			"PolicyDocument": policyDoc.ToMap(),
			"Queues":         []interface{}{queueUrl},
		},
	}, policyID
}

// createTopicPolicy creates an AWS::SNS::TopicPolicy resource.
func (t *ConnectorTransformer) createTopicPolicy(
	logicalID string,
	connector *Connector,
	sourceType, destType, permission string,
	profile *ConnectorProfile,
	templateResources map[string]interface{},
) (map[string]interface{}, string) {
	topicArn := t.getResourceArn(connector.Destination, destType, templateResources)

	// Build policy document
	policyDoc := iam.NewPolicyDocument()
	stmt := iam.NewAllowStatement()
	stmt.Action = "sns:Publish"
	stmt.Resource = topicArn
	stmt.Principal = map[string]interface{}{"Service": "events.amazonaws.com"}
	policyDoc.AddStatement(stmt)

	metadata := t.buildConnectorMetadata(logicalID, sourceType, destType)

	policyID := logicalID + "TopicPolicy"
	return map[string]interface{}{
		"Type":     "AWS::SNS::TopicPolicy",
		"Metadata": metadata,
		"Properties": map[string]interface{}{
			"PolicyDocument": policyDoc.ToMap(),
			"Topics":         []interface{}{topicArn},
		},
	}, policyID
}

// consolidatePolicies merges multiple IAM policies into one if they share the same type.
func (t *ConnectorTransformer) consolidatePolicies(
	logicalID string,
	resources map[string]interface{},
	connector *Connector,
	sourceType, destType string,
) map[string]interface{} {
	// Collect all ManagedPolicy resources
	var policies []map[string]interface{}
	otherResources := make(map[string]interface{})

	for id, r := range resources {
		resource, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if resource["Type"] == "AWS::IAM::ManagedPolicy" {
			policies = append(policies, resource)
		} else {
			otherResources[id] = r
		}
	}

	// If there's only one policy or no policies, no consolidation needed
	if len(policies) <= 1 {
		return resources
	}

	// Merge all statements into one policy document
	mergedDoc := iam.NewPolicyDocument()
	var roles []interface{}

	for _, policy := range policies {
		props, ok := policy["Properties"].(map[string]interface{})
		if !ok {
			continue
		}
		policyDoc, ok := props["PolicyDocument"].(map[string]interface{})
		if !ok {
			continue
		}
		statements, ok := policyDoc["Statement"].([]interface{})
		if !ok {
			continue
		}

		for _, s := range statements {
			stmt, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			newStmt := iam.NewAllowStatement()
			if actions, ok := stmt["Action"]; ok {
				newStmt.Action = actions
			}
			if res, ok := stmt["Resource"]; ok {
				newStmt.Resource = res
			}
			mergedDoc.AddStatement(newStmt)
		}

		if r, ok := props["Roles"].([]interface{}); ok && len(r) > 0 {
			roles = r
		}
	}

	// Build merged policy
	metadata := t.buildConnectorMetadata(logicalID, sourceType, destType)
	mergedPolicy := map[string]interface{}{
		"Type":     "AWS::IAM::ManagedPolicy",
		"Metadata": metadata,
		"Properties": map[string]interface{}{
			"PolicyDocument": mergedDoc.ToMap(),
			"Roles":          roles,
		},
	}

	// Return consolidated resources
	result := otherResources
	result[logicalID+"Policy"] = mergedPolicy
	return result
}

// buildConnectorMetadata builds the aws:sam:connectors metadata.
func (t *ConnectorTransformer) buildConnectorMetadata(logicalID, sourceType, destType string) map[string]interface{} {
	return map[string]interface{}{
		"aws:sam:connectors": map[string]interface{}{
			logicalID: map[string]interface{}{
				"Source": map[string]interface{}{
					"Type": sourceType,
				},
				"Destination": map[string]interface{}{
					"Type": destType,
				},
			},
		},
	}
}

// getRoleReference gets the IAM role reference for a source resource.
func (t *ConnectorTransformer) getRoleReference(endpoint ConnectorEndpoint, resourceType string, templateResources map[string]interface{}) interface{} {
	// If RoleName is explicitly provided, use it
	if endpoint.RoleName != nil {
		return endpoint.RoleName
	}

	// If ID is provided, try to find the role from the resource
	if endpoint.ID != "" {
		resource, ok := templateResources[endpoint.ID]
		if ok {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if props, ok := resourceMap["Properties"].(map[string]interface{}); ok {
					// For Lambda functions, get the Role property and extract role name
					if role, ok := props["Role"]; ok {
						return t.extractRoleNameFromArn(role, endpoint.ID, resourceType)
					}
				}
			}
		}

		// For SAM resources, the role is typically <LogicalID>Role
		if resourceType == TypeServerlessFunction || resourceType == TypeServerlessStateMachine {
			return map[string]interface{}{"Ref": endpoint.ID + "Role"}
		}

		// For Lambda functions with a role property, we need to get the role name
		if resourceType == TypeLambdaFunction {
			if resource, ok := templateResources[endpoint.ID]; ok {
				if resourceMap, ok := resource.(map[string]interface{}); ok {
					if props, ok := resourceMap["Properties"].(map[string]interface{}); ok {
						if role, ok := props["Role"]; ok {
							return t.extractRoleNameFromArn(role, endpoint.ID, resourceType)
						}
					}
				}
			}
		}

		// For Step Functions
		if resourceType == TypeStepFunctionsStateMachine || resourceType == TypeServerlessStateMachine {
			return map[string]interface{}{"Ref": endpoint.ID + "Role"}
		}
	}

	return nil
}

// extractRoleNameFromArn extracts the role name from a role ARN or Fn::GetAtt.
func (t *ConnectorTransformer) extractRoleNameFromArn(role interface{}, resourceID, resourceType string) interface{} {
	// If it's a Fn::GetAtt pointing to a role's Arn, we need the role logical ID
	if roleMap, ok := role.(map[string]interface{}); ok {
		if getAtt, ok := roleMap["Fn::GetAtt"].([]interface{}); ok && len(getAtt) >= 1 {
			roleLogicalID := getAtt[0]
			return map[string]interface{}{"Ref": roleLogicalID}
		}
	}
	// If it's a string ARN, return as is (not ideal but works)
	return role
}

// getResourceArn gets the ARN for a resource endpoint.
func (t *ConnectorTransformer) getResourceArn(endpoint ConnectorEndpoint, resourceType string, templateResources map[string]interface{}) interface{} {
	// If Arn is explicitly provided, use it
	if endpoint.Arn != nil {
		return endpoint.Arn
	}

	// If ID is provided, construct the ARN reference
	if endpoint.ID != "" {
		switch resourceType {
		case TypeSNSTopic:
			return map[string]interface{}{"Ref": endpoint.ID}
		case TypeServerlessStateMachine, TypeStepFunctionsStateMachine:
			return map[string]interface{}{"Ref": endpoint.ID}
		default:
			return map[string]interface{}{
				"Fn::GetAtt": []interface{}{endpoint.ID, "Arn"},
			}
		}
	}

	return nil
}

// getSourceArnForLambdaPermission gets the source ARN for Lambda permission.
func (t *ConnectorTransformer) getSourceArnForLambdaPermission(endpoint ConnectorEndpoint, resourceType string, templateResources map[string]interface{}) interface{} {
	if endpoint.Arn != nil {
		return endpoint.Arn
	}

	if endpoint.ID != "" {
		switch resourceType {
		case TypeSNSTopic:
			return map[string]interface{}{"Ref": endpoint.ID}
		default:
			return map[string]interface{}{
				"Fn::GetAtt": []interface{}{endpoint.ID, "Arn"},
			}
		}
	}

	return nil
}

// getQueueUrl gets the queue URL for an SQS queue.
func (t *ConnectorTransformer) getQueueUrl(endpoint ConnectorEndpoint, templateResources map[string]interface{}) interface{} {
	if endpoint.QueueUrl != nil {
		return endpoint.QueueUrl
	}
	if endpoint.ID != "" {
		return map[string]interface{}{"Ref": endpoint.ID}
	}
	return nil
}

// toInterfaceSlice converts a string slice to an interface slice.
func (t *ConnectorTransformer) toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// ExtractEmbeddedConnectors extracts embedded connectors from a template.
// Returns a map of source resource ID to their embedded connectors.
func ExtractEmbeddedConnectors(templateResources map[string]interface{}) map[string]map[string]EmbeddedConnector {
	result := make(map[string]map[string]EmbeddedConnector)

	for resourceID, resource := range templateResources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			continue
		}

		connectors, ok := resourceMap["Connectors"].(map[string]interface{})
		if !ok {
			continue
		}

		embeddedConnectors := make(map[string]EmbeddedConnector)
		for connectorName, connectorData := range connectors {
			connectorMap, ok := connectorData.(map[string]interface{})
			if !ok {
				continue
			}

			propsData, ok := connectorMap["Properties"].(map[string]interface{})
			if !ok {
				continue
			}

			var dest ConnectorEndpoint
			if destData, ok := propsData["Destination"].(map[string]interface{}); ok {
				if id, ok := destData["Id"].(string); ok {
					dest.ID = id
				}
				if t, ok := destData["Type"].(string); ok {
					dest.Type = t
				}
				if arn, ok := destData["Arn"]; ok {
					dest.Arn = arn
				}
			}

			var permissions []string
			if permsData, ok := propsData["Permissions"].([]interface{}); ok {
				for _, p := range permsData {
					if pstr, ok := p.(string); ok {
						permissions = append(permissions, pstr)
					}
				}
			}

			embeddedConnectors[connectorName] = EmbeddedConnector{
				Properties: EmbeddedConnectorProperties{
					Destination: dest,
					Permissions: permissions,
				},
			}
		}

		if len(embeddedConnectors) > 0 {
			result[resourceID] = embeddedConnectors
		}
	}

	return result
}
