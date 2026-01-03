package intrinsics

import (
	"fmt"
	"regexp"
	"strings"
)

// SubAction handles the Fn::Sub intrinsic function.
// Fn::Sub substitutes variables in a string.
type SubAction struct{}

// Name returns the intrinsic function name.
func (a *SubAction) Name() string {
	return "Fn::Sub"
}

// variablePattern matches ${VarName} or ${VarName.Attribute} patterns
var variablePattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Resolve resolves a Fn::Sub intrinsic function.
// Fn::Sub can take two forms:
// 1. String form: "Fn::Sub": "string with ${variable}"
// 2. Array form: "Fn::Sub": ["string with ${variable}", {"variable": "value"}]
func (a *SubAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return a.substituteString(ctx, v, nil)
	case []interface{}:
		return a.substituteArray(ctx, v)
	default:
		return nil, NewIntrinsicError("Fn::Sub", fmt.Sprintf("expected string or array, got %T", value))
	}
}

// substituteArray handles the array form of Fn::Sub.
func (a *SubAction) substituteArray(ctx *ResolveContext, arr []interface{}) (interface{}, error) {
	if len(arr) != 2 {
		return nil, NewIntrinsicError("Fn::Sub", "array form requires exactly 2 elements")
	}

	templateStr, ok := arr[0].(string)
	if !ok {
		return nil, NewIntrinsicError("Fn::Sub", fmt.Sprintf("first element must be string, got %T", arr[0]))
	}

	varMap, ok := arr[1].(map[string]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::Sub", fmt.Sprintf("second element must be map, got %T", arr[1]))
	}

	return a.substituteString(ctx, templateStr, varMap)
}

// substituteString performs variable substitution on a template string.
func (a *SubAction) substituteString(ctx *ResolveContext, template string, localVars map[string]interface{}) (interface{}, error) {
	// Track unresolved variables for pass-through
	unresolvedVars := make(map[string]bool)
	hasUnresolved := false

	result := variablePattern.ReplaceAllStringFunc(template, func(match string) string {
		// Extract variable name from ${...}
		varName := match[2 : len(match)-1]

		// Check local variables first
		if localVars != nil {
			if val, ok := localVars[varName]; ok {
				return a.valueToString(val)
			}
		}

		// Check for GetAtt syntax (Resource.Attribute)
		if strings.Contains(varName, ".") {
			resolved, err := a.resolveGetAttVar(ctx, varName)
			if err == nil && resolved != nil {
				if str, ok := resolved.(string); ok {
					return str
				}
			}
			// Can't resolve GetAtt - mark as unresolved
			unresolvedVars[varName] = true
			hasUnresolved = true
			return match
		}

		// Check pseudo-parameters
		if pseudo, ok := ctx.PseudoParameters[varName]; ok {
			return pseudo
		}

		// Check parameters
		if param, ok := ctx.Parameters[varName]; ok {
			return a.valueToString(param)
		}

		// Check if it's a resource reference (just the logical ID)
		if ctx.Template != nil && ctx.Template.Resources != nil {
			if _, ok := ctx.Template.Resources[varName]; ok {
				// Resource reference - keep unresolved for CloudFormation
				unresolvedVars[varName] = true
				hasUnresolved = true
				return match
			}
		}

		// Check if it's a defined parameter without a value
		if ctx.Template != nil && ctx.Template.Parameters != nil {
			if _, ok := ctx.Template.Parameters[varName]; ok {
				unresolvedVars[varName] = true
				hasUnresolved = true
				return match
			}
		}

		// Unknown variable - keep for CloudFormation
		unresolvedVars[varName] = true
		hasUnresolved = true
		return match
	})

	// If there are unresolved variables, return the Sub for CloudFormation
	if hasUnresolved {
		if len(localVars) > 0 {
			// Filter localVars to only include unresolved ones
			filteredVars := make(map[string]interface{})
			for k, v := range localVars {
				if unresolvedVars[k] {
					filteredVars[k] = v
				}
			}
			if len(filteredVars) > 0 {
				return map[string]interface{}{
					"Fn::Sub": []interface{}{result, filteredVars},
				}, nil
			}
		}
		return map[string]interface{}{"Fn::Sub": result}, nil
	}

	return result, nil
}

// resolveGetAttVar resolves a GetAtt-style variable reference (Resource.Attribute).
func (a *SubAction) resolveGetAttVar(ctx *ResolveContext, varName string) (interface{}, error) {
	parts := strings.SplitN(varName, ".", 2)
	if len(parts) != 2 {
		return nil, NewIntrinsicError("Fn::Sub", fmt.Sprintf("invalid GetAtt reference: %s", varName))
	}

	resourceName := parts[0]
	attributeName := parts[1]

	// Check if we have cached resource attributes
	if resourceAttrs, ok := ctx.Resources[resourceName]; ok {
		if attrValue, ok := resourceAttrs[attributeName]; ok {
			return attrValue, nil
		}
	}

	// Can't resolve at transform time
	return nil, NewIntrinsicError("Fn::Sub", fmt.Sprintf("cannot resolve %s", varName))
}

// valueToString converts a value to its string representation.
func (a *SubAction) valueToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		// Check if it's a whole number
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
