package intrinsics

import "fmt"

// RefAction handles the Ref intrinsic function.
// Ref returns the value of the specified parameter or resource.
type RefAction struct{}

// Name returns the intrinsic function name.
func (a *RefAction) Name() string {
	return "Ref"
}

// Resolve resolves a Ref intrinsic function.
// For parameters, it returns the parameter value.
// For resources, it returns the resource's physical ID (logical ID for now).
// For pseudo-parameters (AWS::*), it returns the pseudo-parameter value.
// For AWS::NoValue, it returns the NoValue sentinel (use IsNoValue to check).
func (a *RefAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	refName, ok := value.(string)
	if !ok {
		return nil, NewIntrinsicError("Ref", fmt.Sprintf("expected string, got %T", value))
	}

	// AWS::NoValue is special - return sentinel to indicate property removal
	if refName == "AWS::NoValue" {
		return NoValue{}, nil
	}

	// Check pseudo-parameters
	if pseudo, ok := ctx.PseudoParameters[refName]; ok {
		return pseudo, nil
	}

	// Check parameters
	if param, ok := ctx.Parameters[refName]; ok {
		return param, nil
	}

	// Check if it's a defined parameter without a value
	if ctx.Template != nil && ctx.Template.Parameters != nil {
		if _, ok := ctx.Template.Parameters[refName]; ok {
			// Parameter exists but no value - this is an error in real resolution
			// but we keep the Ref for CloudFormation to resolve at deploy time
			return map[string]interface{}{"Ref": refName}, nil
		}
	}

	// Check resources - return the logical ID as a reference
	if ctx.Template != nil && ctx.Template.Resources != nil {
		if _, ok := ctx.Template.Resources[refName]; ok {
			// For resources, we preserve the Ref for CloudFormation
			// since we can't know the physical ID at transform time
			return map[string]interface{}{"Ref": refName}, nil
		}
	}

	// Unknown reference - preserve it for CloudFormation
	return map[string]interface{}{"Ref": refName}, nil
}
