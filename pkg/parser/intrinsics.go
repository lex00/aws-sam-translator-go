// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

// IntrinsicFunctionNames contains all supported CloudFormation intrinsic function names.
var IntrinsicFunctionNames = []string{
	"Ref",
	"Fn::GetAtt",
	"Fn::Sub",
	"Fn::Join",
	"Fn::If",
	"Fn::Select",
	"Fn::FindInMap",
	"Fn::Base64",
	"Fn::Cidr",
	"Fn::GetAZs",
	"Fn::ImportValue",
	"Fn::Split",
	"Fn::Transform",
	"Fn::And",
	"Fn::Equals",
	"Fn::Not",
	"Fn::Or",
	"Condition",
}

// intrinsicFunctionSet provides O(1) lookup for intrinsic function names.
var intrinsicFunctionSet = map[string]bool{
	"Ref":             true,
	"Fn::GetAtt":      true,
	"Fn::Sub":         true,
	"Fn::Join":        true,
	"Fn::If":          true,
	"Fn::Select":      true,
	"Fn::FindInMap":   true,
	"Fn::Base64":      true,
	"Fn::Cidr":        true,
	"Fn::GetAZs":      true,
	"Fn::ImportValue": true,
	"Fn::Split":       true,
	"Fn::Transform":   true,
	"Fn::And":         true,
	"Fn::Equals":      true,
	"Fn::Not":         true,
	"Fn::Or":          true,
	"Condition":       true,
}

// IsIntrinsicFunction checks if a value represents a CloudFormation intrinsic function.
// It returns true if the value is a map with exactly one key that is an intrinsic function name.
func IsIntrinsicFunction(value interface{}) bool {
	m, ok := value.(map[string]interface{})
	if !ok {
		return false
	}
	if len(m) != 1 {
		return false
	}
	for key := range m {
		return intrinsicFunctionSet[key]
	}
	return false
}

// GetIntrinsicName returns the intrinsic function name if the value is an intrinsic,
// otherwise returns an empty string.
func GetIntrinsicName(value interface{}) string {
	m, ok := value.(map[string]interface{})
	if !ok {
		return ""
	}
	if len(m) != 1 {
		return ""
	}
	for key := range m {
		if intrinsicFunctionSet[key] {
			return key
		}
	}
	return ""
}

// GetIntrinsicValue returns the value of an intrinsic function.
// Returns nil if the value is not an intrinsic function.
func GetIntrinsicValue(value interface{}) interface{} {
	m, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	for key, val := range m {
		if intrinsicFunctionSet[key] {
			return val
		}
	}
	return nil
}

// IntrinsicInfo contains information about a detected intrinsic function.
type IntrinsicInfo struct {
	Name     string
	Value    interface{}
	Path     string
	Location SourceLocation
	Nested   []IntrinsicInfo
}

// DetectIntrinsics recursively scans a value and returns all intrinsic functions found.
func DetectIntrinsics(value interface{}, path string) []IntrinsicInfo {
	var intrinsics []IntrinsicInfo
	detectIntrinsicsRecursive(value, path, &intrinsics)
	return intrinsics
}

// detectIntrinsicsRecursive is the recursive helper for DetectIntrinsics.
func detectIntrinsicsRecursive(value interface{}, path string, intrinsics *[]IntrinsicInfo) {
	switch v := value.(type) {
	case map[string]interface{}:
		if name := GetIntrinsicName(v); name != "" {
			info := IntrinsicInfo{
				Name:  name,
				Value: v[name],
				Path:  path,
			}
			// Check for nested intrinsics within this intrinsic's value
			info.Nested = DetectIntrinsics(v[name], path+"."+name)
			*intrinsics = append(*intrinsics, info)
		} else {
			// Not an intrinsic, check children
			for key, val := range v {
				childPath := path
				if childPath != "" {
					childPath += "." + key
				} else {
					childPath = key
				}
				detectIntrinsicsRecursive(val, childPath, intrinsics)
			}
		}
	case []interface{}:
		for i, item := range v {
			itemPath := path + "[" + intToString(i) + "]"
			detectIntrinsicsRecursive(item, itemPath, intrinsics)
		}
	}
}

// ContainsIntrinsics checks if a value contains any intrinsic functions.
func ContainsIntrinsics(value interface{}) bool {
	switch v := value.(type) {
	case map[string]interface{}:
		if IsIntrinsicFunction(v) {
			return true
		}
		for _, val := range v {
			if ContainsIntrinsics(val) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if ContainsIntrinsics(item) {
				return true
			}
		}
	}
	return false
}

// CountIntrinsics returns the number of intrinsic functions in a value.
func CountIntrinsics(value interface{}) int {
	return len(DetectIntrinsics(value, ""))
}

// ValidateIntrinsicStructure validates that an intrinsic function has the correct structure.
// Returns an error if the structure is invalid.
func ValidateIntrinsicStructure(name string, value interface{}) error {
	switch name {
	case "Ref":
		return validateRef(value)
	case "Fn::GetAtt":
		return validateGetAtt(value)
	case "Fn::Sub":
		return validateSub(value)
	case "Fn::Join":
		return validateJoin(value)
	case "Fn::If":
		return validateIf(value)
	case "Fn::Select":
		return validateSelect(value)
	case "Fn::FindInMap":
		return validateFindInMap(value)
	case "Fn::Base64":
		return validateBase64(value)
	case "Fn::Cidr":
		return validateCidr(value)
	case "Fn::GetAZs":
		return validateGetAZs(value)
	case "Fn::And", "Fn::Or":
		return validateAndOr(value)
	case "Fn::Equals":
		return validateEquals(value)
	case "Fn::Not":
		return validateNot(value)
	default:
		return nil
	}
}

func validateRef(value interface{}) error {
	if _, ok := value.(string); !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Ref",
			Message:   "Ref value must be a string",
		}
	}
	return nil
}

func validateGetAtt(value interface{}) error {
	switch v := value.(type) {
	case string:
		// String form: "Resource.Attribute"
		return nil
	case []interface{}:
		if len(v) != 2 {
			return &IntrinsicValidationError{
				Intrinsic: "Fn::GetAtt",
				Message:   "Fn::GetAtt array must have exactly 2 elements",
			}
		}
	default:
		return &IntrinsicValidationError{
			Intrinsic: "Fn::GetAtt",
			Message:   "Fn::GetAtt value must be a string or array",
		}
	}
	return nil
}

func validateSub(value interface{}) error {
	switch v := value.(type) {
	case string:
		return nil
	case []interface{}:
		if len(v) != 2 {
			return &IntrinsicValidationError{
				Intrinsic: "Fn::Sub",
				Message:   "Fn::Sub array must have exactly 2 elements",
			}
		}
		if _, ok := v[0].(string); !ok {
			return &IntrinsicValidationError{
				Intrinsic: "Fn::Sub",
				Message:   "Fn::Sub first element must be a string template",
			}
		}
		if _, ok := v[1].(map[string]interface{}); !ok {
			return &IntrinsicValidationError{
				Intrinsic: "Fn::Sub",
				Message:   "Fn::Sub second element must be a map of variables",
			}
		}
	default:
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Sub",
			Message:   "Fn::Sub value must be a string or array",
		}
	}
	return nil
}

func validateJoin(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Join",
			Message:   "Fn::Join value must be an array",
		}
	}
	if len(arr) != 2 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Join",
			Message:   "Fn::Join array must have exactly 2 elements",
		}
	}
	if _, ok := arr[0].(string); !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Join",
			Message:   "Fn::Join first element must be a delimiter string",
		}
	}
	if _, ok := arr[1].([]interface{}); !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Join",
			Message:   "Fn::Join second element must be an array",
		}
	}
	return nil
}

func validateIf(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::If",
			Message:   "Fn::If value must be an array",
		}
	}
	if len(arr) != 3 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::If",
			Message:   "Fn::If array must have exactly 3 elements",
		}
	}
	if _, ok := arr[0].(string); !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::If",
			Message:   "Fn::If first element must be a condition name string",
		}
	}
	return nil
}

func validateSelect(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Select",
			Message:   "Fn::Select value must be an array",
		}
	}
	if len(arr) != 2 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Select",
			Message:   "Fn::Select array must have exactly 2 elements",
		}
	}
	return nil
}

func validateFindInMap(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::FindInMap",
			Message:   "Fn::FindInMap value must be an array",
		}
	}
	if len(arr) != 3 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::FindInMap",
			Message:   "Fn::FindInMap array must have exactly 3 elements",
		}
	}
	return nil
}

func validateBase64(value interface{}) error {
	switch value.(type) {
	case string, map[string]interface{}:
		return nil
	default:
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Base64",
			Message:   "Fn::Base64 value must be a string or intrinsic function",
		}
	}
}

func validateCidr(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Cidr",
			Message:   "Fn::Cidr value must be an array",
		}
	}
	if len(arr) != 3 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Cidr",
			Message:   "Fn::Cidr array must have exactly 3 elements",
		}
	}
	return nil
}

func validateGetAZs(value interface{}) error {
	switch value.(type) {
	case string, map[string]interface{}:
		return nil
	default:
		return &IntrinsicValidationError{
			Intrinsic: "Fn::GetAZs",
			Message:   "Fn::GetAZs value must be a string or intrinsic function",
		}
	}
}

func validateAndOr(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::And/Fn::Or",
			Message:   "Fn::And/Fn::Or value must be an array",
		}
	}
	if len(arr) < 2 || len(arr) > 10 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::And/Fn::Or",
			Message:   "Fn::And/Fn::Or array must have 2-10 elements",
		}
	}
	return nil
}

func validateEquals(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Equals",
			Message:   "Fn::Equals value must be an array",
		}
	}
	if len(arr) != 2 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Equals",
			Message:   "Fn::Equals array must have exactly 2 elements",
		}
	}
	return nil
}

func validateNot(value interface{}) error {
	arr, ok := value.([]interface{})
	if !ok {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Not",
			Message:   "Fn::Not value must be an array",
		}
	}
	if len(arr) != 1 {
		return &IntrinsicValidationError{
			Intrinsic: "Fn::Not",
			Message:   "Fn::Not array must have exactly 1 element",
		}
	}
	return nil
}

// IntrinsicValidationError represents an error in intrinsic function structure.
type IntrinsicValidationError struct {
	Intrinsic string
	Message   string
	Location  SourceLocation
}

func (e *IntrinsicValidationError) Error() string {
	if e.Location.Line > 0 {
		return e.Intrinsic + " at line " + intToString(e.Location.Line) + ": " + e.Message
	}
	return e.Intrinsic + ": " + e.Message
}
