package intrinsics

import "fmt"

// JoinAction handles the Fn::Join intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type JoinAction struct{}

// Name returns the intrinsic function name.
func (a *JoinAction) Name() string {
	return "Fn::Join"
}

// Resolve preserves the Fn::Join intrinsic for CloudFormation.
// While we could evaluate joins with static values, we preserve them
// for CloudFormation to maintain template fidelity.
func (a *JoinAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::Join", fmt.Sprintf("expected array, got %T", value))
	}
	if len(arr) != 2 {
		return nil, NewIntrinsicError("Fn::Join", fmt.Sprintf("expected 2 elements, got %d", len(arr)))
	}
	return map[string]interface{}{"Fn::Join": value}, nil
}

// IfAction handles the Fn::If intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type IfAction struct{}

// Name returns the intrinsic function name.
func (a *IfAction) Name() string {
	return "Fn::If"
}

// Resolve processes Fn::If, potentially evaluating it if the condition is known.
func (a *IfAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::If", fmt.Sprintf("expected array, got %T", value))
	}
	if len(arr) != 3 {
		return nil, NewIntrinsicError("Fn::If", fmt.Sprintf("expected 3 elements, got %d", len(arr)))
	}

	conditionName, ok := arr[0].(string)
	if !ok {
		return nil, NewIntrinsicError("Fn::If", fmt.Sprintf("condition name must be string, got %T", arr[0]))
	}

	// Check if we know the condition value
	if condValue, ok := ctx.Conditions[conditionName]; ok {
		if condValue {
			return arr[1], nil // Return true value
		}
		return arr[2], nil // Return false value
	}

	// Condition not evaluated - preserve for CloudFormation
	return map[string]interface{}{"Fn::If": value}, nil
}

// SelectAction handles the Fn::Select intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type SelectAction struct{}

// Name returns the intrinsic function name.
func (a *SelectAction) Name() string {
	return "Fn::Select"
}

// Resolve preserves the Fn::Select intrinsic for CloudFormation.
func (a *SelectAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::Select", fmt.Sprintf("expected array, got %T", value))
	}
	if len(arr) != 2 {
		return nil, NewIntrinsicError("Fn::Select", fmt.Sprintf("expected 2 elements, got %d", len(arr)))
	}
	return map[string]interface{}{"Fn::Select": value}, nil
}

// Base64Action handles the Fn::Base64 intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type Base64Action struct{}

// Name returns the intrinsic function name.
func (a *Base64Action) Name() string {
	return "Fn::Base64"
}

// Resolve preserves the Fn::Base64 intrinsic for CloudFormation.
func (a *Base64Action) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	return map[string]interface{}{"Fn::Base64": value}, nil
}

// GetAZsAction handles the Fn::GetAZs intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type GetAZsAction struct{}

// Name returns the intrinsic function name.
func (a *GetAZsAction) Name() string {
	return "Fn::GetAZs"
}

// Resolve preserves the Fn::GetAZs intrinsic for CloudFormation.
func (a *GetAZsAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	return map[string]interface{}{"Fn::GetAZs": value}, nil
}

// SplitAction handles the Fn::Split intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type SplitAction struct{}

// Name returns the intrinsic function name.
func (a *SplitAction) Name() string {
	return "Fn::Split"
}

// Resolve preserves the Fn::Split intrinsic for CloudFormation.
func (a *SplitAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::Split", fmt.Sprintf("expected array, got %T", value))
	}
	if len(arr) != 2 {
		return nil, NewIntrinsicError("Fn::Split", fmt.Sprintf("expected 2 elements, got %d", len(arr)))
	}
	return map[string]interface{}{"Fn::Split": value}, nil
}

// ImportValueAction handles the Fn::ImportValue intrinsic function.
// This is a pass-through handler that preserves the intrinsic for CloudFormation.
type ImportValueAction struct{}

// Name returns the intrinsic function name.
func (a *ImportValueAction) Name() string {
	return "Fn::ImportValue"
}

// Resolve preserves the Fn::ImportValue intrinsic for CloudFormation.
func (a *ImportValueAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	return map[string]interface{}{"Fn::ImportValue": value}, nil
}

// ConditionAction handles the Condition intrinsic function.
// This is used to reference a condition in resource properties.
type ConditionAction struct{}

// Name returns the intrinsic function name.
func (a *ConditionAction) Name() string {
	return "Condition"
}

// Resolve evaluates a condition reference if the condition is known.
func (a *ConditionAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	conditionName, ok := value.(string)
	if !ok {
		return nil, NewIntrinsicError("Condition", fmt.Sprintf("expected string, got %T", value))
	}

	// Check if we know the condition value
	if condValue, ok := ctx.Conditions[conditionName]; ok {
		return condValue, nil
	}

	// Condition not evaluated - preserve for CloudFormation
	return map[string]interface{}{"Condition": value}, nil
}
