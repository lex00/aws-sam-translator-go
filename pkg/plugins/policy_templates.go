package plugins

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/policy"
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// PolicyTemplatesPlugin expands SAM policy templates in function Policies.
type PolicyTemplatesPlugin struct {
	processor *policy.Processor
}

// NewPolicyTemplatesPlugin creates a new PolicyTemplatesPlugin.
func NewPolicyTemplatesPlugin() (*PolicyTemplatesPlugin, error) {
	processor, err := policy.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create policy processor: %w", err)
	}

	return &PolicyTemplatesPlugin{
		processor: processor,
	}, nil
}

// Name returns the plugin name.
func (p *PolicyTemplatesPlugin) Name() string {
	return "PolicyTemplatesPlugin"
}

// Priority returns the execution priority (400).
func (p *PolicyTemplatesPlugin) Priority() int {
	return 400
}

// BeforeTransform expands policy templates in all functions.
func (p *PolicyTemplatesPlugin) BeforeTransform(template *types.Template) error {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::Function" {
			if err := p.expandFunctionPolicies(&resource); err != nil {
				return err
			}
		}
	}

	return nil
}

// AfterTransform does nothing for this plugin.
func (p *PolicyTemplatesPlugin) AfterTransform(template *types.Template) error {
	return nil
}

// expandFunctionPolicies expands policy templates in a function's Policies property.
func (p *PolicyTemplatesPlugin) expandFunctionPolicies(resource *types.Resource) error {
	if resource.Properties == nil {
		return nil
	}

	policies, ok := resource.Properties["Policies"]
	if !ok {
		return nil
	}

	// Policies can be a string, array, or object
	switch pol := policies.(type) {
	case string:
		// Single managed policy ARN - nothing to expand
		return nil

	case []interface{}:
		// Array of policies - may contain template references
		expanded, err := p.expandPolicyArray(pol)
		if err != nil {
			return err
		}
		resource.Properties["Policies"] = expanded

	case map[string]interface{}:
		// Single policy object - may be a template
		expanded, err := p.expandPolicyObject(pol)
		if err != nil {
			return err
		}
		resource.Properties["Policies"] = []interface{}{expanded}

	default:
		// Unknown type - leave as-is
		return nil
	}

	return nil
}

// expandPolicyArray expands policy templates in an array of policies.
func (p *PolicyTemplatesPlugin) expandPolicyArray(policies []interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0, len(policies))

	for _, pol := range policies {
		switch policy := pol.(type) {
		case string:
			// Managed policy ARN - keep as-is
			result = append(result, policy)

		case map[string]interface{}:
			expanded, err := p.expandPolicyObject(policy)
			if err != nil {
				return nil, err
			}
			result = append(result, expanded)

		default:
			// Unknown type - keep as-is
			result = append(result, pol)
		}
	}

	return result, nil
}

// expandPolicyObject expands a single policy object if it's a SAM policy template.
func (p *PolicyTemplatesPlugin) expandPolicyObject(policy map[string]interface{}) (map[string]interface{}, error) {
	// Check if this is a SAM policy template (has a single key that matches a template name)
	if len(policy) == 1 {
		for templateName, params := range policy {
			// Check if this is a known policy template
			if p.processor.HasTemplate(templateName) {
				// Extract parameters
				paramMap, ok := params.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("policy template %s parameters must be a map", templateName)
				}

				// Expand the template
				expanded, err := p.processor.Expand(templateName, paramMap)
				if err != nil {
					return nil, fmt.Errorf("failed to expand policy template %s: %w", templateName, err)
				}

				return expanded, nil
			}
		}
	}

	// Not a SAM policy template - return as-is
	return policy, nil
}
