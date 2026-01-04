// Package translator provides the main SAM to CloudFormation transformation orchestrator.
package translator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/parser"
	"github.com/lex00/aws-sam-translator-go/pkg/plugins"
	"github.com/lex00/aws-sam-translator-go/pkg/sam"
	"github.com/lex00/aws-sam-translator-go/pkg/types"
	"github.com/lex00/cloudformation-schema-go/spec"
)

// Version is the translator library version.
const Version = "0.1.0"

// SAMTransform is the SAM transform identifier.
const SAMTransform = "AWS::Serverless-2016-10-31"

// Options configures the translator behavior.
type Options struct {
	// Region is the AWS region for resource ARN generation.
	Region string

	// AccountID is the AWS account ID for resource ARN generation.
	AccountID string

	// StackName is the CloudFormation stack name.
	StackName string

	// Partition is the AWS partition (aws, aws-cn, aws-us-gov).
	Partition string

	// PassThroughMetadata preserves template-level Metadata in output.
	PassThroughMetadata bool

	// FeatureToggles controls optional transformation features.
	FeatureToggles map[string]bool
}

// Translator transforms SAM templates to CloudFormation.
type Translator struct {
	schema         *spec.Spec
	options        Options
	pluginRegistry *plugins.Registry

	// Transformers for each SAM resource type
	functionTransformer     *sam.FunctionTransformer
	simpleTableTransformer  *sam.SimpleTableTransformer
	layerVersionTransformer *sam.LayerVersionTransformer
	stateMachineTransformer *sam.StateMachineTransformer
	apiTransformer          *sam.ApiTransformer
	httpApiTransformer      *sam.HttpApiTransformer
	applicationTransformer  *sam.ApplicationTransformer
	graphQLApiTransformer   *sam.GraphQLApiTransformer
	connectorTransformer    *sam.ConnectorTransformer
}

// Schema returns the CloudFormation schema.
func (t *Translator) Schema() *spec.Spec {
	return t.schema
}

// New creates a new Translator instance with default options.
func New() *Translator {
	return NewWithOptions(Options{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: "aws",
	})
}

// NewWithOptions creates a new Translator instance with the specified options.
func NewWithOptions(opts Options) *Translator {
	// Apply defaults for empty values
	if opts.Region == "" {
		opts.Region = "us-east-1"
	}
	if opts.AccountID == "" {
		opts.AccountID = "123456789012"
	}
	if opts.StackName == "" {
		opts.StackName = "sam-app"
	}
	if opts.Partition == "" {
		opts.Partition = "aws"
	}

	t := &Translator{
		options:                 opts,
		pluginRegistry:          plugins.NewRegistry(),
		functionTransformer:     sam.NewFunctionTransformer(),
		simpleTableTransformer:  sam.NewSimpleTableTransformer(),
		layerVersionTransformer: sam.NewLayerVersionTransformer(),
		stateMachineTransformer: sam.NewStateMachineTransformer(),
		apiTransformer:          sam.NewApiTransformer(),
		httpApiTransformer:      sam.NewHttpApiTransformer(),
		applicationTransformer:  sam.NewApplicationTransformer(),
		graphQLApiTransformer:   sam.NewGraphQLApiTransformer(),
		connectorTransformer:    sam.NewConnectorTransformer(),
	}

	// Register default plugins
	t.registerDefaultPlugins()

	return t
}

// registerDefaultPlugins registers the built-in SAM plugins.
func (t *Translator) registerDefaultPlugins() {
	t.pluginRegistry.Register(plugins.NewGlobalsPlugin())
	// PolicyTemplatesPlugin returns error, but we ignore it for now as the plugin
	// loader will work with embedded templates
	if policyPlugin, err := plugins.NewPolicyTemplatesPlugin(); err == nil {
		t.pluginRegistry.Register(policyPlugin)
	}
	t.pluginRegistry.Register(plugins.NewImplicitRestApiPlugin())
	t.pluginRegistry.Register(plugins.NewImplicitHttpApiPlugin())
	t.pluginRegistry.Register(plugins.NewDefaultDefinitionBodyPlugin())
}

// RegisterPlugin registers an additional plugin.
func (t *Translator) RegisterPlugin(p plugins.Plugin) {
	t.pluginRegistry.Register(p)
}

// Transform converts a SAM template to CloudFormation.
func (t *Translator) Transform(template *types.Template) (*types.Template, error) {
	// Create the output template
	output := &types.Template{
		AWSTemplateFormatVersion: template.AWSTemplateFormatVersion,
		Description:              template.Description,
		Parameters:               template.Parameters,
		Mappings:                 template.Mappings,
		Conditions:               template.Conditions,
		Outputs:                  template.Outputs,
		Resources:                make(map[string]types.Resource),
	}

	// Handle metadata passthrough
	if t.options.PassThroughMetadata && template.Metadata != nil {
		output.Metadata = template.Metadata
	}

	// Handle Transform - remove SAM transform, preserve others
	output.Transform = t.filterTransform(template.Transform)

	// Run BeforeTransform plugins
	if err := t.pluginRegistry.RunBeforeTransform(template); err != nil {
		return nil, fmt.Errorf("BeforeTransform plugin error: %w", err)
	}

	// Create transform context
	ctx := &sam.TransformContext{
		Region:    t.options.Region,
		AccountID: t.options.AccountID,
		StackName: t.options.StackName,
		Partition: t.options.Partition,
	}

	// Get ordered list of resources to transform
	orderedResources := t.getOrderedResources(template.Resources)

	// Track all errors for aggregation
	var errs []error

	// Transform each resource in order
	for _, entry := range orderedResources {
		logicalID := entry.logicalID
		resource := entry.resource

		if isSAMResource(resource.Type) {
			// Transform SAM resource
			newResources, err := t.transformSAMResource(logicalID, resource, ctx, template)
			if err != nil {
				errs = append(errs, fmt.Errorf("resource '%s': %w", logicalID, err))
				continue
			}

			// Add transformed resources to output
			for id, res := range newResources {
				output.Resources[id] = res
			}
		} else {
			// Pass through non-SAM resources unchanged
			output.Resources[logicalID] = resource
		}
	}

	// Run AfterTransform plugins
	if err := t.pluginRegistry.RunAfterTransform(output); err != nil {
		return nil, fmt.Errorf("AfterTransform plugin error: %w", err)
	}

	// Return aggregated errors if any
	if len(errs) > 0 {
		return nil, &TransformError{Errors: errs}
	}

	return output, nil
}

// TransformBytes parses a YAML/JSON template and transforms it to CloudFormation JSON.
func (t *Translator) TransformBytes(input []byte) ([]byte, error) {
	// Parse the input template
	p := parser.New()
	template, err := p.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Transform
	result, err := t.Transform(template)
	if err != nil {
		return nil, err
	}

	// Marshal to JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output: %w", err)
	}

	return output, nil
}

// resourceEntry holds a resource with its logical ID for sorting.
type resourceEntry struct {
	logicalID string
	resource  types.Resource
	priority  int
}

// getOrderedResources returns resources in the correct processing order.
func (t *Translator) getOrderedResources(resources map[string]types.Resource) []resourceEntry {
	entries := make([]resourceEntry, 0, len(resources))

	order := getResourceOrder()
	priorityMap := make(map[string]int)
	for i, rt := range order {
		priorityMap[rt] = i
	}

	for logicalID, resource := range resources {
		priority := 999 // Default priority for unknown types
		if p, ok := priorityMap[resource.Type]; ok {
			priority = p
		}
		entries = append(entries, resourceEntry{
			logicalID: logicalID,
			resource:  resource,
			priority:  priority,
		})
	}

	// Sort by priority, then by logical ID for determinism
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].priority != entries[j].priority {
			return entries[i].priority < entries[j].priority
		}
		return entries[i].logicalID < entries[j].logicalID
	})

	return entries
}

// getResourceOrder returns the order in which SAM resource types should be processed.
// Order matters because:
// - Functions must be processed before APIs (APIs may reference functions)
// - SimpleTable/LayerVersion are independent
// - StateMachines may reference functions
// - Connectors must be processed last (they reference other resources)
func getResourceOrder() []string {
	return []string{
		"AWS::Serverless::LayerVersion", // Independent, process first
		"AWS::Serverless::SimpleTable",  // Independent
		"AWS::Serverless::Function",     // Functions before APIs
		"AWS::Serverless::StateMachine", // May reference functions
		"AWS::Serverless::Api",          // APIs after functions
		"AWS::Serverless::HttpApi",      // HTTP APIs after functions
		"AWS::Serverless::Application",  // Nested stacks
		"AWS::Serverless::GraphQLApi",   // AppSync APIs
		"AWS::Serverless::Connector",    // Last - references other resources
	}
}

// isSAMResource checks if a resource type is a SAM resource.
func isSAMResource(resourceType string) bool {
	return strings.HasPrefix(resourceType, "AWS::Serverless::")
}

// transformSAMResource transforms a single SAM resource.
func (t *Translator) transformSAMResource(logicalID string, resource types.Resource, ctx *sam.TransformContext, template *types.Template) (map[string]types.Resource, error) {
	switch resource.Type {
	case "AWS::Serverless::Function":
		return t.transformFunction(logicalID, resource, ctx)
	case "AWS::Serverless::SimpleTable":
		return t.transformSimpleTable(logicalID, resource)
	case "AWS::Serverless::LayerVersion":
		return t.transformLayerVersion(logicalID, resource)
	case "AWS::Serverless::StateMachine":
		return t.transformStateMachine(logicalID, resource, ctx)
	case "AWS::Serverless::Api":
		return t.transformApi(logicalID, resource, ctx)
	case "AWS::Serverless::HttpApi":
		return t.transformHttpApi(logicalID, resource, ctx)
	case "AWS::Serverless::Application":
		return t.transformApplication(logicalID, resource, ctx)
	case "AWS::Serverless::GraphQLApi":
		return t.transformGraphQLApi(logicalID, resource, ctx)
	case "AWS::Serverless::Connector":
		return t.transformConnector(logicalID, resource, template)
	default:
		return nil, fmt.Errorf("unknown SAM resource type: %s", resource.Type)
	}
}

// transformFunction transforms an AWS::Serverless::Function resource.
func (t *Translator) transformFunction(logicalID string, resource types.Resource, ctx *sam.TransformContext) (map[string]types.Resource, error) {
	fn, err := t.parseFunction(resource.Properties)
	if err != nil {
		return nil, err
	}

	// Copy resource-level properties
	fn.Condition = resource.Condition
	fn.DependsOn = resource.DependsOn
	fn.Metadata = resource.Metadata

	rawResources, err := t.functionTransformer.Transform(logicalID, fn, ctx)
	if err != nil {
		return nil, err
	}

	return t.convertRawResources(rawResources), nil
}

// transformSimpleTable transforms an AWS::Serverless::SimpleTable resource.
func (t *Translator) transformSimpleTable(logicalID string, resource types.Resource) (map[string]types.Resource, error) {
	st, err := t.parseSimpleTable(resource.Properties)
	if err != nil {
		return nil, err
	}

	rawResources, err := t.simpleTableTransformer.Transform(logicalID, st)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties
	if r, ok := result[logicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[logicalID] = r
	}

	return result, nil
}

// transformLayerVersion transforms an AWS::Serverless::LayerVersion resource.
func (t *Translator) transformLayerVersion(logicalID string, resource types.Resource) (map[string]types.Resource, error) {
	lv, err := t.parseLayerVersion(resource.Properties)
	if err != nil {
		return nil, err
	}

	rawResources, newLogicalID, err := t.layerVersionTransformer.Transform(logicalID, lv)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties to the new logical ID
	if r, ok := result[newLogicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[newLogicalID] = r
	}

	return result, nil
}

// transformStateMachine transforms an AWS::Serverless::StateMachine resource.
func (t *Translator) transformStateMachine(logicalID string, resource types.Resource, _ *sam.TransformContext) (map[string]types.Resource, error) {
	sm, err := t.parseStateMachine(resource.Properties)
	if err != nil {
		return nil, err
	}

	rawResources, err := t.stateMachineTransformer.Transform(logicalID, sm)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties to the main state machine resource
	if r, ok := result[logicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[logicalID] = r
	}

	return result, nil
}

// transformApi transforms an AWS::Serverless::Api resource.
func (t *Translator) transformApi(logicalID string, resource types.Resource, _ *sam.TransformContext) (map[string]types.Resource, error) {
	api, err := t.parseApi(resource.Properties)
	if err != nil {
		return nil, err
	}

	rawResources, err := t.apiTransformer.Transform(logicalID, api)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties to the main Api resource
	if r, ok := result[logicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[logicalID] = r
	}

	return result, nil
}

// transformHttpApi transforms an AWS::Serverless::HttpApi resource.
func (t *Translator) transformHttpApi(logicalID string, resource types.Resource, ctx *sam.TransformContext) (map[string]types.Resource, error) {
	httpApi, err := t.parseHttpApi(resource.Properties)
	if err != nil {
		return nil, err
	}

	rawResources, err := t.httpApiTransformer.Transform(logicalID, httpApi, ctx)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties to the main HttpApi resource
	if r, ok := result[logicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[logicalID] = r
	}

	return result, nil
}

// transformApplication transforms an AWS::Serverless::Application resource.
func (t *Translator) transformApplication(logicalID string, resource types.Resource, ctx *sam.TransformContext) (map[string]types.Resource, error) {
	app, err := t.parseApplication(resource.Properties)
	if err != nil {
		return nil, err
	}

	// Copy resource-level properties (Application has these fields)
	app.Condition = resource.Condition
	app.DependsOn = resource.DependsOn
	app.Metadata = resource.Metadata

	rawResources, err := t.applicationTransformer.Transform(logicalID, app, ctx)
	if err != nil {
		return nil, err
	}

	return t.convertRawResources(rawResources), nil
}

// transformGraphQLApi transforms an AWS::Serverless::GraphQLApi resource.
func (t *Translator) transformGraphQLApi(logicalID string, resource types.Resource, ctx *sam.TransformContext) (map[string]types.Resource, error) {
	gql, err := t.parseGraphQLApi(resource.Properties)
	if err != nil {
		return nil, err
	}

	// Copy resource-level properties
	gql.Condition = resource.Condition
	gql.DependsOn = resource.DependsOn
	gql.Metadata = resource.Metadata

	rawResources, err := t.graphQLApiTransformer.Transform(logicalID, gql, ctx)
	if err != nil {
		return nil, err
	}

	return t.convertRawResources(rawResources), nil
}

// transformConnector transforms an AWS::Serverless::Connector resource.
func (t *Translator) transformConnector(logicalID string, resource types.Resource, template *types.Template) (map[string]types.Resource, error) {
	conn, err := t.parseConnector(resource.Properties)
	if err != nil {
		return nil, err
	}

	// Convert template resources to map[string]interface{} for connector transformer
	templateResources := make(map[string]interface{})
	for id, res := range template.Resources {
		templateResources[id] = map[string]interface{}{
			"Type":       res.Type,
			"Properties": res.Properties,
		}
	}

	rawResources, err := t.connectorTransformer.Transform(logicalID, conn, templateResources)
	if err != nil {
		return nil, err
	}

	result := t.convertRawResources(rawResources)

	// Copy resource-level properties to main connector resource if present
	if r, ok := result[logicalID]; ok {
		r.Condition = resource.Condition
		r.DependsOn = resource.DependsOn
		r.Metadata = resource.Metadata
		result[logicalID] = r
	}

	return result, nil
}

// convertRawResources converts a map[string]interface{} to map[string]types.Resource.
func (t *Translator) convertRawResources(raw map[string]interface{}) map[string]types.Resource {
	result := make(map[string]types.Resource)

	for logicalID, rawRes := range raw {
		resMap, ok := rawRes.(map[string]interface{})
		if !ok {
			continue
		}

		resource := types.Resource{}

		if resType, ok := resMap["Type"].(string); ok {
			resource.Type = resType
		}
		if props, ok := resMap["Properties"].(map[string]interface{}); ok {
			resource.Properties = props
		}
		if metadata, ok := resMap["Metadata"].(map[string]interface{}); ok {
			resource.Metadata = metadata
		}
		if condition, ok := resMap["Condition"].(string); ok {
			resource.Condition = condition
		}
		if dependsOn, ok := resMap["DependsOn"]; ok {
			resource.DependsOn = dependsOn
		}
		if deletionPolicy, ok := resMap["DeletionPolicy"].(string); ok {
			resource.DeletionPolicy = deletionPolicy
		}
		if updatePolicy, ok := resMap["UpdatePolicy"].(map[string]interface{}); ok {
			resource.UpdatePolicy = updatePolicy
		}

		result[logicalID] = resource
	}

	return result
}

// filterTransform removes the SAM transform from the Transform value.
func (t *Translator) filterTransform(transform interface{}) interface{} {
	if transform == nil {
		return nil
	}

	switch v := transform.(type) {
	case string:
		if v == SAMTransform {
			return nil
		}
		return v
	case []interface{}:
		var filtered []interface{}
		for _, item := range v {
			if s, ok := item.(string); ok && s == SAMTransform {
				continue
			}
			filtered = append(filtered, item)
		}
		if len(filtered) == 0 {
			return nil
		}
		if len(filtered) == 1 {
			return filtered[0]
		}
		return filtered
	default:
		return transform
	}
}

// TransformError aggregates multiple transformation errors.
type TransformError struct {
	Errors []error
}

func (e *TransformError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("multiple errors:\n  - %s", strings.Join(msgs, "\n  - "))
}
