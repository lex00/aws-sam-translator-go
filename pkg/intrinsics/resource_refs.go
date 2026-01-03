// Package intrinsics provides CloudFormation intrinsic function handling.
package intrinsics

import (
	"fmt"
	"sort"
)

// DependencyTracker tracks resource dependencies discovered during intrinsic resolution.
// It builds a dependency graph based on Ref and GetAtt intrinsic function usage.
type DependencyTracker struct {
	// dependencies maps a resource to the resources it depends on
	dependencies map[string]map[string]bool
	// dependents maps a resource to the resources that depend on it
	dependents map[string]map[string]bool
}

// NewDependencyTracker creates a new DependencyTracker.
func NewDependencyTracker() *DependencyTracker {
	return &DependencyTracker{
		dependencies: make(map[string]map[string]bool),
		dependents:   make(map[string]map[string]bool),
	}
}

// AddDependency records that sourceResource depends on targetResource.
func (dt *DependencyTracker) AddDependency(sourceResource, targetResource string) {
	// Initialize maps if needed
	if dt.dependencies[sourceResource] == nil {
		dt.dependencies[sourceResource] = make(map[string]bool)
	}
	if dt.dependents[targetResource] == nil {
		dt.dependents[targetResource] = make(map[string]bool)
	}

	// Record the dependency
	dt.dependencies[sourceResource][targetResource] = true
	dt.dependents[targetResource][sourceResource] = true
}

// GetDependencies returns the resources that a given resource depends on.
func (dt *DependencyTracker) GetDependencies(resourceName string) []string {
	deps := dt.dependencies[resourceName]
	if deps == nil {
		return nil
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	sort.Strings(result)
	return result
}

// GetDependents returns the resources that depend on a given resource.
func (dt *DependencyTracker) GetDependents(resourceName string) []string {
	deps := dt.dependents[resourceName]
	if deps == nil {
		return nil
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	sort.Strings(result)
	return result
}

// HasDependency checks if sourceResource depends on targetResource.
func (dt *DependencyTracker) HasDependency(sourceResource, targetResource string) bool {
	if deps := dt.dependencies[sourceResource]; deps != nil {
		return deps[targetResource]
	}
	return false
}

// AllDependencies returns the complete dependency map.
func (dt *DependencyTracker) AllDependencies() map[string][]string {
	result := make(map[string][]string)
	for resource, deps := range dt.dependencies {
		depList := make([]string, 0, len(deps))
		for dep := range deps {
			depList = append(depList, dep)
		}
		sort.Strings(depList)
		result[resource] = depList
	}
	return result
}

// TopologicalSort returns resources in dependency order (dependencies first).
// Returns an error if there's a circular dependency.
func (dt *DependencyTracker) TopologicalSort(resources []string) ([]string, error) {
	// Build a set of resources to sort
	resourceSet := make(map[string]bool)
	for _, r := range resources {
		resourceSet[r] = true
	}

	// Track visited and visiting nodes for cycle detection
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	result := make([]string, 0, len(resources))

	// Helper function for DFS
	var visit func(resource string) error
	visit = func(resource string) error {
		if visited[resource] {
			return nil
		}
		if visiting[resource] {
			return &CircularDependencyError{Resource: resource}
		}

		visiting[resource] = true

		// Visit all dependencies first
		for dep := range dt.dependencies[resource] {
			if resourceSet[dep] {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		visiting[resource] = false
		visited[resource] = true
		result = append(result, resource)
		return nil
	}

	// Sort input resources for deterministic output
	sortedResources := make([]string, len(resources))
	copy(sortedResources, resources)
	sort.Strings(sortedResources)

	// Visit each resource
	for _, resource := range sortedResources {
		if err := visit(resource); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// CircularDependencyError indicates a circular dependency was detected.
type CircularDependencyError struct {
	Resource string
}

func (e *CircularDependencyError) Error() string {
	return fmt.Sprintf("circular dependency detected involving resource: %s", e.Resource)
}

// ResourceRef represents a reference to a resource for tracking purposes.
type ResourceRef struct {
	// LogicalID is the logical ID of the referenced resource
	LogicalID string
	// Attribute is the attribute being referenced (empty for Ref, populated for GetAtt)
	Attribute string
	// SourcePath is where this reference was found
	SourcePath string
}

// ResourceRefCollector collects all resource references in a template.
type ResourceRefCollector struct {
	refs     []ResourceRef
	template interface{}
}

// NewResourceRefCollector creates a new ResourceRefCollector.
func NewResourceRefCollector() *ResourceRefCollector {
	return &ResourceRefCollector{
		refs: make([]ResourceRef, 0),
	}
}

// CollectRefs traverses a value and collects all resource references.
func (c *ResourceRefCollector) CollectRefs(value interface{}, path string) []ResourceRef {
	c.refs = make([]ResourceRef, 0)
	c.collectRefsFromValue(value, path)
	return c.refs
}

// collectRefsFromValue recursively collects references from a value.
func (c *ResourceRefCollector) collectRefsFromValue(value interface{}, path string) {
	switch v := value.(type) {
	case map[string]interface{}:
		c.collectRefsFromMap(v, path)
	case []interface{}:
		for i, elem := range v {
			c.collectRefsFromValue(elem, fmt.Sprintf("%s[%d]", path, i))
		}
	}
}

// collectRefsFromMap handles map values, looking for Ref and GetAtt.
func (c *ResourceRefCollector) collectRefsFromMap(m map[string]interface{}, path string) {
	// Check for Ref
	if refVal, ok := m["Ref"]; ok {
		if len(m) == 1 {
			if refName, ok := refVal.(string); ok {
				c.refs = append(c.refs, ResourceRef{
					LogicalID:  refName,
					SourcePath: path,
				})
			}
		}
	}

	// Check for Fn::GetAtt
	if getAttVal, ok := m["Fn::GetAtt"]; ok {
		if len(m) == 1 {
			switch v := getAttVal.(type) {
			case []interface{}:
				if len(v) >= 2 {
					if logicalID, ok := v[0].(string); ok {
						if attr, ok := v[1].(string); ok {
							c.refs = append(c.refs, ResourceRef{
								LogicalID:  logicalID,
								Attribute:  attr,
								SourcePath: path,
							})
						}
					}
				}
			case string:
				// Short form: "LogicalID.Attribute"
				parts := splitFirst(v, ".")
				if len(parts) == 2 {
					c.refs = append(c.refs, ResourceRef{
						LogicalID:  parts[0],
						Attribute:  parts[1],
						SourcePath: path,
					})
				}
			}
		}
	}

	// Recurse into all values
	for key, val := range m {
		c.collectRefsFromValue(val, path+"."+key)
	}
}

// splitFirst splits a string on the first occurrence of sep.
func splitFirst(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}

// MergeDependencies merges another DependencyTracker into this one.
func (dt *DependencyTracker) MergeDependencies(other *DependencyTracker) {
	if other == nil {
		return
	}

	for source, deps := range other.dependencies {
		for target := range deps {
			dt.AddDependency(source, target)
		}
	}
}

// Clear resets the dependency tracker.
func (dt *DependencyTracker) Clear() {
	dt.dependencies = make(map[string]map[string]bool)
	dt.dependents = make(map[string]map[string]bool)
}

// Count returns the total number of dependency relationships tracked.
func (dt *DependencyTracker) Count() int {
	count := 0
	for _, deps := range dt.dependencies {
		count += len(deps)
	}
	return count
}
