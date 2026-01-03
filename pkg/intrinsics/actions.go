// Package intrinsics provides CloudFormation intrinsic function handling.
package intrinsics

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// ResolveContext provides context for resolving intrinsic functions.
type ResolveContext struct {
	// Template is the full template being processed
	Template *types.Template
	// Parameters contains parameter values (defaults or provided)
	Parameters map[string]interface{}
	// Resources tracks resource logical IDs and their resolved properties
	Resources map[string]map[string]interface{}
	// Conditions tracks condition evaluations (true/false)
	Conditions map[string]bool
	// PseudoParameters contains AWS pseudo-parameter values
	PseudoParameters map[string]string
}

// NewResolveContext creates a new ResolveContext with initialized maps.
func NewResolveContext(template *types.Template) *ResolveContext {
	ctx := &ResolveContext{
		Template:         template,
		Parameters:       make(map[string]interface{}),
		Resources:        make(map[string]map[string]interface{}),
		Conditions:       make(map[string]bool),
		PseudoParameters: make(map[string]string),
	}

	// Initialize default pseudo-parameters
	ctx.PseudoParameters["AWS::AccountId"] = "123456789012"
	ctx.PseudoParameters["AWS::Region"] = "us-east-1"
	ctx.PseudoParameters["AWS::StackName"] = "sam-app"
	ctx.PseudoParameters["AWS::StackId"] = "arn:aws:cloudformation:us-east-1:123456789012:stack/sam-app/guid"
	ctx.PseudoParameters["AWS::URLSuffix"] = "amazonaws.com"
	ctx.PseudoParameters["AWS::Partition"] = "aws"
	ctx.PseudoParameters["AWS::NoValue"] = ""

	// Extract parameter defaults from template
	if template != nil && template.Parameters != nil {
		for name, param := range template.Parameters {
			if param.Default != nil {
				ctx.Parameters[name] = param.Default
			}
		}
	}

	return ctx
}

// SetPseudoParameter sets a pseudo-parameter value.
func (ctx *ResolveContext) SetPseudoParameter(name, value string) {
	ctx.PseudoParameters[name] = value
}

// SetParameter sets a parameter value.
func (ctx *ResolveContext) SetParameter(name string, value interface{}) {
	ctx.Parameters[name] = value
}

// Action defines the interface for intrinsic function handlers.
type Action interface {
	// Name returns the intrinsic function name (e.g., "Ref", "Fn::Sub")
	Name() string
	// Resolve processes the intrinsic function and returns the resolved value.
	// The value parameter contains the intrinsic's argument.
	// Returns the resolved value or an error if resolution fails.
	Resolve(ctx *ResolveContext, value interface{}) (interface{}, error)
}

// Registry holds registered intrinsic function actions.
type Registry struct {
	actions map[string]Action
}

// NewRegistry creates a new action registry with default actions.
func NewRegistry() *Registry {
	r := &Registry{
		actions: make(map[string]Action),
	}
	r.registerDefaults()
	return r
}

// registerDefaults registers all built-in intrinsic function handlers.
func (r *Registry) registerDefaults() {
	r.Register(&RefAction{})
	r.Register(&SubAction{})
	r.Register(&GetAttAction{})
	r.Register(&FindInMapAction{})
	r.Register(&JoinAction{})
	r.Register(&IfAction{})
	r.Register(&SelectAction{})
	r.Register(&Base64Action{})
	r.Register(&GetAZsAction{})
	r.Register(&SplitAction{})
	r.Register(&ImportValueAction{})
	r.Register(&ConditionAction{})
}

// Register adds an action to the registry.
func (r *Registry) Register(action Action) {
	r.actions[action.Name()] = action
}

// Get returns the action for a given intrinsic function name.
func (r *Registry) Get(name string) (Action, bool) {
	action, ok := r.actions[name]
	return action, ok
}

// Resolve resolves all intrinsic functions in a value recursively.
func (r *Registry) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	return r.resolveValue(ctx, value)
}

// resolveValue recursively resolves intrinsic functions in a value.
func (r *Registry) resolveValue(ctx *ResolveContext, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return r.resolveMap(ctx, v)
	case []interface{}:
		return r.resolveSlice(ctx, v)
	default:
		return value, nil
	}
}

// resolveMap resolves intrinsic functions in a map.
func (r *Registry) resolveMap(ctx *ResolveContext, m map[string]interface{}) (interface{}, error) {
	// Check if this map is an intrinsic function
	if len(m) == 1 {
		for key, val := range m {
			if action, ok := r.actions[key]; ok {
				// First resolve any nested intrinsics in the value
				resolvedVal, err := r.resolveValue(ctx, val)
				if err != nil {
					return nil, err
				}
				// Then apply this intrinsic's action
				return action.Resolve(ctx, resolvedVal)
			}
		}
	}

	// Not an intrinsic, resolve all values in the map
	result := make(map[string]interface{})
	for key, val := range m {
		resolved, err := r.resolveValue(ctx, val)
		if err != nil {
			return nil, err
		}
		result[key] = resolved
	}
	return result, nil
}

// resolveSlice resolves intrinsic functions in a slice.
func (r *Registry) resolveSlice(ctx *ResolveContext, s []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(s))
	for i, val := range s {
		resolved, err := r.resolveValue(ctx, val)
		if err != nil {
			return nil, err
		}
		result[i] = resolved
	}
	return result, nil
}

// IntrinsicError represents an error during intrinsic function resolution.
type IntrinsicError struct {
	Function string
	Message  string
}

func (e *IntrinsicError) Error() string {
	return fmt.Sprintf("%s: %s", e.Function, e.Message)
}

// NewIntrinsicError creates a new IntrinsicError.
func NewIntrinsicError(function, message string) *IntrinsicError {
	return &IntrinsicError{Function: function, Message: message}
}
