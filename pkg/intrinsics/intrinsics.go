// Package intrinsics provides CloudFormation intrinsic function handling.
package intrinsics

// Ref represents a CloudFormation Ref intrinsic.
type Ref struct {
	Ref string `json:"Ref" yaml:"Ref"`
}

// GetAtt represents a CloudFormation Fn::GetAtt intrinsic.
type GetAtt struct {
	GetAtt []string `json:"Fn::GetAtt" yaml:"Fn::GetAtt"`
}

// Sub represents a CloudFormation Fn::Sub intrinsic.
type Sub struct {
	Sub interface{} `json:"Fn::Sub" yaml:"Fn::Sub"`
}

// Join represents a CloudFormation Fn::Join intrinsic.
type Join struct {
	Join []interface{} `json:"Fn::Join" yaml:"Fn::Join"`
}

// If represents a CloudFormation Fn::If intrinsic.
type If struct {
	If []interface{} `json:"Fn::If" yaml:"Fn::If"`
}

// Select represents a CloudFormation Fn::Select intrinsic.
type Select struct {
	Select []interface{} `json:"Fn::Select" yaml:"Fn::Select"`
}

// FindInMap represents a CloudFormation Fn::FindInMap intrinsic.
type FindInMap struct {
	FindInMap []interface{} `json:"Fn::FindInMap" yaml:"Fn::FindInMap"`
}

// IsIntrinsic checks if a value is a CloudFormation intrinsic function.
func IsIntrinsic(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	if len(m) != 1 {
		return false
	}
	for k := range m {
		switch k {
		case "Ref", "Fn::GetAtt", "Fn::Sub", "Fn::Join", "Fn::If", "Fn::Select", "Fn::FindInMap",
			"Fn::Base64", "Fn::Cidr", "Fn::GetAZs", "Fn::ImportValue", "Fn::Split", "Fn::Transform",
			"Fn::And", "Fn::Equals", "Fn::Not", "Fn::Or", "Condition":
			return true
		}
	}
	return false
}
