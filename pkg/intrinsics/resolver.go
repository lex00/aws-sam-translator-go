// Package intrinsics provides CloudFormation intrinsic function handling.
package intrinsics

import (
	"fmt"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// Resolver handles intrinsic function resolution across a template with support
// for pre-order tree traversal, logical ID mutation, and placeholder protection.
type Resolver struct {
	registry      *Registry
	context       *ResolveContext
	dependencies  *DependencyTracker
	logicalIDMap  map[string]string // Maps old logical IDs to new logical IDs
	placeholders  map[string]interface{}
}

// ResolverOption configures the Resolver.
type ResolverOption func(*Resolver)

// WithLogicalIDMap sets a mapping from old to new logical IDs for mutation support.
func WithLogicalIDMap(idMap map[string]string) ResolverOption {
	return func(r *Resolver) {
		r.logicalIDMap = idMap
	}
}

// WithPlaceholders sets values that should be protected from resolution.
// These are typically values that need CloudFormation runtime resolution.
func WithPlaceholders(placeholders map[string]interface{}) ResolverOption {
	return func(r *Resolver) {
		r.placeholders = placeholders
	}
}

// WithDependencyTracker sets a custom dependency tracker.
func WithDependencyTracker(tracker *DependencyTracker) ResolverOption {
	return func(r *Resolver) {
		r.dependencies = tracker
	}
}

// NewResolver creates a new Resolver with the given context and options.
func NewResolver(ctx *ResolveContext, opts ...ResolverOption) *Resolver {
	r := &Resolver{
		registry:     NewRegistry(),
		context:      ctx,
		dependencies: NewDependencyTracker(),
		logicalIDMap: make(map[string]string),
		placeholders: make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Resolve processes a value using pre-order tree traversal, resolving all
// intrinsic functions while respecting logical ID mutations and placeholders.
func (r *Resolver) Resolve(value interface{}) (interface{}, error) {
	return r.resolveValue(value, "")
}

// ResolveTemplate processes an entire template, resolving intrinsics in all sections.
func (r *Resolver) ResolveTemplate(template *types.Template) (*types.Template, error) {
	if template == nil {
		return nil, nil
	}

	// Create a copy of the template to avoid mutating the original
	result := &types.Template{
		AWSTemplateFormatVersion: template.AWSTemplateFormatVersion,
		Transform:                template.Transform,
		Description:              template.Description,
	}

	// Resolve Metadata
	if template.Metadata != nil {
		resolved, err := r.resolveMap(template.Metadata, "Metadata")
		if err != nil {
			return nil, fmt.Errorf("resolving Metadata: %w", err)
		}
		result.Metadata = resolved
	}

	// Copy Parameters (parameters are definitions, not resolved)
	result.Parameters = template.Parameters

	// Copy Mappings (mappings are static lookup tables)
	result.Mappings = template.Mappings

	// Resolve Conditions
	if template.Conditions != nil {
		resolved, err := r.resolveMap(template.Conditions, "Conditions")
		if err != nil {
			return nil, fmt.Errorf("resolving Conditions: %w", err)
		}
		result.Conditions = resolved
	}

	// Resolve Resources
	if template.Resources != nil {
		result.Resources = make(map[string]types.Resource)
		for logicalID, resource := range template.Resources {
			// Apply logical ID mutation if configured
			newLogicalID := r.mapLogicalID(logicalID)

			resolvedResource, err := r.resolveResource(resource, newLogicalID)
			if err != nil {
				return nil, fmt.Errorf("resolving resource %s: %w", logicalID, err)
			}
			result.Resources[newLogicalID] = resolvedResource
		}
	}

	// Resolve Outputs
	if template.Outputs != nil {
		result.Outputs = make(map[string]types.Output)
		for name, output := range template.Outputs {
			resolvedOutput, err := r.resolveOutput(output, name)
			if err != nil {
				return nil, fmt.Errorf("resolving output %s: %w", name, err)
			}
			result.Outputs[name] = resolvedOutput
		}
	}

	// Copy Globals (SAM-specific, typically handled separately)
	result.Globals = template.Globals

	return result, nil
}

// resolveResource resolves intrinsics in a resource's properties and metadata.
func (r *Resolver) resolveResource(resource types.Resource, logicalID string) (types.Resource, error) {
	result := types.Resource{
		Type:           resource.Type,
		Condition:      resource.Condition,
		DeletionPolicy: resource.DeletionPolicy,
	}

	// Resolve Properties
	if resource.Properties != nil {
		resolved, err := r.resolveMap(resource.Properties, fmt.Sprintf("Resources.%s.Properties", logicalID))
		if err != nil {
			return result, err
		}
		result.Properties = resolved
	}

	// Resolve Metadata
	if resource.Metadata != nil {
		resolved, err := r.resolveMap(resource.Metadata, fmt.Sprintf("Resources.%s.Metadata", logicalID))
		if err != nil {
			return result, err
		}
		result.Metadata = resolved
	}

	// Resolve UpdatePolicy
	if resource.UpdatePolicy != nil {
		resolved, err := r.resolveMap(resource.UpdatePolicy, fmt.Sprintf("Resources.%s.UpdatePolicy", logicalID))
		if err != nil {
			return result, err
		}
		result.UpdatePolicy = resolved
	}

	// Handle DependsOn with logical ID mapping
	if resource.DependsOn != nil {
		result.DependsOn = r.mapDependsOn(resource.DependsOn)
	}

	return result, nil
}

// resolveOutput resolves intrinsics in an output's value and export.
func (r *Resolver) resolveOutput(output types.Output, name string) (types.Output, error) {
	result := types.Output{
		Description: output.Description,
		Condition:   output.Condition,
	}

	// Resolve Value
	if output.Value != nil {
		resolved, err := r.resolveValue(output.Value, fmt.Sprintf("Outputs.%s.Value", name))
		if err != nil {
			return result, err
		}
		result.Value = resolved
	}

	// Resolve Export
	if output.Export != nil {
		result.Export = &types.Export{}
		if output.Export.Name != nil {
			resolved, err := r.resolveValue(output.Export.Name, fmt.Sprintf("Outputs.%s.Export.Name", name))
			if err != nil {
				return result, err
			}
			result.Export.Name = resolved
		}
	}

	return result, nil
}

// resolveValue performs pre-order tree traversal resolution on a value.
// Pre-order means we process the current node (intrinsic) before its children,
// but the Registry already handles the nested resolution in the correct order.
func (r *Resolver) resolveValue(value interface{}, path string) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return r.resolveMapValue(v, path)
	case []interface{}:
		return r.resolveSlice(v, path)
	default:
		return value, nil
	}
}

// resolveMapValue handles map values, detecting and resolving intrinsics.
func (r *Resolver) resolveMapValue(m map[string]interface{}, path string) (interface{}, error) {
	// Check if this is an intrinsic function (single-key map with intrinsic name)
	if len(m) == 1 {
		for key, val := range m {
			if r.isIntrinsicKey(key) {
				// Apply logical ID mutation to the value before resolution
				mutatedVal := r.mutateLogicalIDs(val)

				// Check for placeholders that should be preserved
				if r.isPlaceholder(key, mutatedVal) {
					return m, nil
				}

				// Track dependencies for Ref and GetAtt
				r.trackDependency(key, mutatedVal, path)

				// Resolve nested values first (pre-order: process current, then recurse)
				resolvedVal, err := r.resolveValue(mutatedVal, path+"."+key)
				if err != nil {
					return nil, err
				}

				// Apply the intrinsic action
				action, ok := r.registry.Get(key)
				if !ok {
					// Unknown intrinsic - preserve for CloudFormation
					return map[string]interface{}{key: resolvedVal}, nil
				}

				return action.Resolve(r.context, resolvedVal)
			}
		}
	}

	// Not an intrinsic - resolve all nested values
	return r.resolveMap(m, path)
}

// resolveMap resolves all values in a map.
func (r *Resolver) resolveMap(m map[string]interface{}, path string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, val := range m {
		resolved, err := r.resolveValue(val, path+"."+key)
		if err != nil {
			return nil, err
		}

		// Handle NoValue - remove the property
		if IsNoValue(resolved) {
			continue
		}

		result[key] = resolved
	}
	return result, nil
}

// resolveSlice resolves all values in a slice.
func (r *Resolver) resolveSlice(s []interface{}, path string) ([]interface{}, error) {
	result := make([]interface{}, 0, len(s))
	for i, val := range s {
		resolved, err := r.resolveValue(val, fmt.Sprintf("%s[%d]", path, i))
		if err != nil {
			return nil, err
		}

		// Handle NoValue - skip the element
		if IsNoValue(resolved) {
			continue
		}

		result = append(result, resolved)
	}
	return result, nil
}

// isIntrinsicKey checks if a key represents an intrinsic function.
func (r *Resolver) isIntrinsicKey(key string) bool {
	switch key {
	case "Ref", "Condition":
		return true
	default:
		return strings.HasPrefix(key, "Fn::")
	}
}

// mapLogicalID returns the new logical ID for a given old logical ID.
func (r *Resolver) mapLogicalID(logicalID string) string {
	if newID, ok := r.logicalIDMap[logicalID]; ok {
		return newID
	}
	return logicalID
}

// mutateLogicalIDs applies logical ID mapping to values.
func (r *Resolver) mutateLogicalIDs(value interface{}) interface{} {
	if len(r.logicalIDMap) == 0 {
		return value
	}

	switch v := value.(type) {
	case string:
		// For Ref values - direct logical ID reference
		if newID, ok := r.logicalIDMap[v]; ok {
			return newID
		}
		return v
	case []interface{}:
		// For GetAtt values - [LogicalID, Attribute]
		if len(v) >= 1 {
			if logicalID, ok := v[0].(string); ok {
				if newID, ok := r.logicalIDMap[logicalID]; ok {
					result := make([]interface{}, len(v))
					result[0] = newID
					copy(result[1:], v[1:])
					return result
				}
			}
		}
		return v
	case map[string]interface{}:
		// Recursively mutate nested structures
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = r.mutateLogicalIDs(val)
		}
		return result
	default:
		return value
	}
}

// mapDependsOn applies logical ID mapping to DependsOn values.
func (r *Resolver) mapDependsOn(dependsOn interface{}) interface{} {
	switch v := dependsOn.(type) {
	case string:
		return r.mapLogicalID(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, dep := range v {
			if depStr, ok := dep.(string); ok {
				result[i] = r.mapLogicalID(depStr)
			} else {
				result[i] = dep
			}
		}
		return result
	case []string:
		result := make([]string, len(v))
		for i, dep := range v {
			result[i] = r.mapLogicalID(dep)
		}
		return result
	default:
		return dependsOn
	}
}

// isPlaceholder checks if a value should be protected from resolution.
func (r *Resolver) isPlaceholder(intrinsicKey string, value interface{}) bool {
	// Check explicit placeholders
	if strVal, ok := value.(string); ok {
		if _, isPlaceholder := r.placeholders[strVal]; isPlaceholder {
			return true
		}
	}

	// Certain intrinsics should always be preserved for CloudFormation runtime
	switch intrinsicKey {
	case "Fn::ImportValue", "Fn::GetAZs":
		return true
	}

	return false
}

// trackDependency records resource dependencies from Ref and GetAtt.
func (r *Resolver) trackDependency(intrinsicKey string, value interface{}, path string) {
	// Extract the source resource from the path
	sourceResource := r.extractResourceFromPath(path)
	if sourceResource == "" {
		return
	}

	switch intrinsicKey {
	case "Ref":
		if refName, ok := value.(string); ok {
			// Only track resource references, not parameters or pseudo-parameters
			if r.isResourceReference(refName) {
				r.dependencies.AddDependency(sourceResource, refName)
			}
		}
	case "Fn::GetAtt":
		var targetResource string
		switch v := value.(type) {
		case []interface{}:
			if len(v) >= 1 {
				if res, ok := v[0].(string); ok {
					targetResource = res
				}
			}
		case string:
			// Short form: "ResourceName.AttributeName"
			parts := strings.SplitN(v, ".", 2)
			if len(parts) >= 1 {
				targetResource = parts[0]
			}
		}
		if targetResource != "" {
			r.dependencies.AddDependency(sourceResource, targetResource)
		}
	}
}

// extractResourceFromPath extracts the resource logical ID from a path like "Resources.MyFunc.Properties".
func (r *Resolver) extractResourceFromPath(path string) string {
	if !strings.HasPrefix(path, "Resources.") {
		return ""
	}
	parts := strings.SplitN(path, ".", 3)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// isResourceReference checks if a reference name refers to a resource.
func (r *Resolver) isResourceReference(name string) bool {
	// Check if it's a pseudo-parameter
	if strings.HasPrefix(name, "AWS::") {
		return false
	}

	// Check if it's a parameter
	if _, ok := r.context.Parameters[name]; ok {
		return false
	}

	// Check if it's in the template parameters
	if r.context.Template != nil && r.context.Template.Parameters != nil {
		if _, ok := r.context.Template.Parameters[name]; ok {
			return false
		}
	}

	// Assume it's a resource reference
	return true
}

// GetDependencies returns the dependency tracker for inspection.
func (r *Resolver) GetDependencies() *DependencyTracker {
	return r.dependencies
}

// SetParameter sets a parameter value in the resolver's context.
func (r *Resolver) SetParameter(name string, value interface{}) {
	r.context.SetParameter(name, value)
}

// SetResourceAttribute sets a resource attribute value for GetAtt resolution.
func (r *Resolver) SetResourceAttribute(resourceName, attributeName string, value interface{}) {
	if r.context.Resources[resourceName] == nil {
		r.context.Resources[resourceName] = make(map[string]interface{})
	}
	r.context.Resources[resourceName][attributeName] = value
}

// AddLogicalIDMapping adds a mapping from an old logical ID to a new one.
func (r *Resolver) AddLogicalIDMapping(oldID, newID string) {
	r.logicalIDMap[oldID] = newID
}

// AddPlaceholder marks a value as protected from resolution.
func (r *Resolver) AddPlaceholder(name string, value interface{}) {
	r.placeholders[name] = value
}
