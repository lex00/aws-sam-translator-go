package intrinsics

import (
	"fmt"
	"strings"
)

// GetAttAction handles the Fn::GetAtt intrinsic function.
// Fn::GetAtt returns the value of an attribute from a resource.
type GetAttAction struct{}

// Name returns the intrinsic function name.
func (a *GetAttAction) Name() string {
	return "Fn::GetAtt"
}

// Resolve resolves a Fn::GetAtt intrinsic function.
// Fn::GetAtt can take two forms:
// 1. Array form: ["logicalName", "attributeName"]
// 2. String form (short syntax): "logicalName.attributeName"
func (a *GetAttAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	var resourceName, attributeName string

	switch v := value.(type) {
	case []interface{}:
		if len(v) < 2 {
			return nil, NewIntrinsicError("Fn::GetAtt", "array must have at least 2 elements")
		}
		var ok bool
		resourceName, ok = v[0].(string)
		if !ok {
			return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("resource name must be string, got %T", v[0]))
		}
		attributeName, ok = v[1].(string)
		if !ok {
			return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("attribute name must be string, got %T", v[1]))
		}
		// Handle nested attributes (e.g., ["Resource", "Attr1", "Attr2"])
		if len(v) > 2 {
			parts := make([]string, len(v)-1)
			for i := 1; i < len(v); i++ {
				part, ok := v[i].(string)
				if !ok {
					return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("attribute part must be string, got %T", v[i]))
				}
				parts[i-1] = part
			}
			attributeName = strings.Join(parts, ".")
		}
	case string:
		// Short form: "ResourceName.AttributeName"
		parts := strings.SplitN(v, ".", 2)
		if len(parts) != 2 {
			return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("invalid string format: %s", v))
		}
		resourceName = parts[0]
		attributeName = parts[1]
	default:
		return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("expected array or string, got %T", value))
	}

	// Check if we have the resource in our context
	if ctx.Template != nil && ctx.Template.Resources != nil {
		if _, ok := ctx.Template.Resources[resourceName]; !ok {
			return nil, NewIntrinsicError("Fn::GetAtt", fmt.Sprintf("resource '%s' not found", resourceName))
		}
	}

	// Check if we have cached resource attributes
	if resourceAttrs, ok := ctx.Resources[resourceName]; ok {
		if attrValue, ok := resourceAttrs[attributeName]; ok {
			return attrValue, nil
		}
	}

	// Most GetAtt lookups need to be resolved by CloudFormation at deploy time
	// We preserve the intrinsic for CloudFormation to handle
	return map[string]interface{}{
		"Fn::GetAtt": []interface{}{resourceName, attributeName},
	}, nil
}
