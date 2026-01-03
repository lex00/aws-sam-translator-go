package iam

import (
	"fmt"
)

// Policy represents an AWS::IAM::Policy CloudFormation resource.
type Policy struct {
	// PolicyName is the name of the policy.
	PolicyName interface{} `json:"PolicyName"`

	// PolicyDocument is the policy document.
	PolicyDocument *PolicyDocument `json:"PolicyDocument"`

	// Groups is a list of IAM group names to attach the policy to.
	Groups []interface{} `json:"Groups,omitempty"`

	// Roles is a list of IAM role names to attach the policy to.
	Roles []interface{} `json:"Roles,omitempty"`

	// Users is a list of IAM user names to attach the policy to.
	Users []interface{} `json:"Users,omitempty"`
}

// NewPolicy creates a new Policy with the specified name and document.
func NewPolicy(name interface{}, document *PolicyDocument) *Policy {
	return &Policy{
		PolicyName:     name,
		PolicyDocument: document,
	}
}

// AttachToRole attaches the policy to an IAM role.
func (p *Policy) AttachToRole(roleName interface{}) *Policy {
	p.Roles = append(p.Roles, roleName)
	return p
}

// AttachToRoles attaches the policy to multiple IAM roles.
func (p *Policy) AttachToRoles(roleNames []interface{}) *Policy {
	p.Roles = append(p.Roles, roleNames...)
	return p
}

// AttachToUser attaches the policy to an IAM user.
func (p *Policy) AttachToUser(userName interface{}) *Policy {
	p.Users = append(p.Users, userName)
	return p
}

// AttachToUsers attaches the policy to multiple IAM users.
func (p *Policy) AttachToUsers(userNames []interface{}) *Policy {
	p.Users = append(p.Users, userNames...)
	return p
}

// AttachToGroup attaches the policy to an IAM group.
func (p *Policy) AttachToGroup(groupName interface{}) *Policy {
	p.Groups = append(p.Groups, groupName)
	return p
}

// AttachToGroups attaches the policy to multiple IAM groups.
func (p *Policy) AttachToGroups(groupNames []interface{}) *Policy {
	p.Groups = append(p.Groups, groupNames...)
	return p
}

// ToCloudFormation converts the policy to CloudFormation resource properties.
func (p *Policy) ToCloudFormation() map[string]interface{} {
	props := make(map[string]interface{})

	props["PolicyName"] = p.PolicyName

	if p.PolicyDocument != nil {
		props["PolicyDocument"] = p.PolicyDocument.ToMap()
	}

	if len(p.Groups) > 0 {
		props["Groups"] = p.Groups
	}

	if len(p.Roles) > 0 {
		props["Roles"] = p.Roles
	}

	if len(p.Users) > 0 {
		props["Users"] = p.Users
	}

	return props
}

// ToResource converts the policy to a complete CloudFormation resource.
func (p *Policy) ToResource() map[string]interface{} {
	return map[string]interface{}{
		"Type":       "AWS::IAM::Policy",
		"Properties": p.ToCloudFormation(),
	}
}

// Validate validates the policy configuration.
func (p *Policy) Validate() error {
	if p.PolicyName == nil {
		return fmt.Errorf("PolicyName is required for IAM::Policy")
	}

	if p.PolicyDocument == nil {
		return fmt.Errorf("PolicyDocument is required for IAM::Policy")
	}

	if err := p.PolicyDocument.Validate(); err != nil {
		return fmt.Errorf("invalid PolicyDocument: %w", err)
	}

	// At least one attachment target is required
	if len(p.Groups) == 0 && len(p.Roles) == 0 && len(p.Users) == 0 {
		return fmt.Errorf("IAM::Policy must be attached to at least one Group, Role, or User")
	}

	return nil
}

// ManagedPolicy represents an AWS::IAM::ManagedPolicy CloudFormation resource.
type ManagedPolicy struct {
	// ManagedPolicyName is the name of the managed policy.
	ManagedPolicyName interface{} `json:"ManagedPolicyName,omitempty"`

	// PolicyDocument is the policy document.
	PolicyDocument *PolicyDocument `json:"PolicyDocument"`

	// Description is a description of the managed policy.
	Description string `json:"Description,omitempty"`

	// Path is the path for the managed policy.
	Path string `json:"Path,omitempty"`

	// Groups is a list of IAM group names to attach the policy to.
	Groups []interface{} `json:"Groups,omitempty"`

	// Roles is a list of IAM role names to attach the policy to.
	Roles []interface{} `json:"Roles,omitempty"`

	// Users is a list of IAM user names to attach the policy to.
	Users []interface{} `json:"Users,omitempty"`
}

// NewManagedPolicy creates a new ManagedPolicy with the specified document.
func NewManagedPolicy(document *PolicyDocument) *ManagedPolicy {
	return &ManagedPolicy{
		PolicyDocument: document,
	}
}

// NewManagedPolicyWithName creates a new ManagedPolicy with a name and document.
func NewManagedPolicyWithName(name interface{}, document *PolicyDocument) *ManagedPolicy {
	return &ManagedPolicy{
		ManagedPolicyName: name,
		PolicyDocument:    document,
	}
}

// WithDescription sets the description for the managed policy.
func (p *ManagedPolicy) WithDescription(description string) *ManagedPolicy {
	p.Description = description
	return p
}

// WithPath sets the path for the managed policy.
func (p *ManagedPolicy) WithPath(path string) *ManagedPolicy {
	p.Path = path
	return p
}

// AttachToRole attaches the managed policy to an IAM role.
func (p *ManagedPolicy) AttachToRole(roleName interface{}) *ManagedPolicy {
	p.Roles = append(p.Roles, roleName)
	return p
}

// AttachToRoles attaches the managed policy to multiple IAM roles.
func (p *ManagedPolicy) AttachToRoles(roleNames []interface{}) *ManagedPolicy {
	p.Roles = append(p.Roles, roleNames...)
	return p
}

// AttachToUser attaches the managed policy to an IAM user.
func (p *ManagedPolicy) AttachToUser(userName interface{}) *ManagedPolicy {
	p.Users = append(p.Users, userName)
	return p
}

// AttachToUsers attaches the managed policy to multiple IAM users.
func (p *ManagedPolicy) AttachToUsers(userNames []interface{}) *ManagedPolicy {
	p.Users = append(p.Users, userNames...)
	return p
}

// AttachToGroup attaches the managed policy to an IAM group.
func (p *ManagedPolicy) AttachToGroup(groupName interface{}) *ManagedPolicy {
	p.Groups = append(p.Groups, groupName)
	return p
}

// AttachToGroups attaches the managed policy to multiple IAM groups.
func (p *ManagedPolicy) AttachToGroups(groupNames []interface{}) *ManagedPolicy {
	p.Groups = append(p.Groups, groupNames...)
	return p
}

// ToCloudFormation converts the managed policy to CloudFormation resource properties.
func (p *ManagedPolicy) ToCloudFormation() map[string]interface{} {
	props := make(map[string]interface{})

	if p.ManagedPolicyName != nil {
		props["ManagedPolicyName"] = p.ManagedPolicyName
	}

	if p.PolicyDocument != nil {
		props["PolicyDocument"] = p.PolicyDocument.ToMap()
	}

	if p.Description != "" {
		props["Description"] = p.Description
	}

	if p.Path != "" {
		props["Path"] = p.Path
	}

	if len(p.Groups) > 0 {
		props["Groups"] = p.Groups
	}

	if len(p.Roles) > 0 {
		props["Roles"] = p.Roles
	}

	if len(p.Users) > 0 {
		props["Users"] = p.Users
	}

	return props
}

// ToResource converts the managed policy to a complete CloudFormation resource.
func (p *ManagedPolicy) ToResource() map[string]interface{} {
	return map[string]interface{}{
		"Type":       "AWS::IAM::ManagedPolicy",
		"Properties": p.ToCloudFormation(),
	}
}

// Validate validates the managed policy configuration.
func (p *ManagedPolicy) Validate() error {
	if p.PolicyDocument == nil {
		return fmt.Errorf("PolicyDocument is required for IAM::ManagedPolicy")
	}

	if err := p.PolicyDocument.Validate(); err != nil {
		return fmt.Errorf("invalid PolicyDocument: %w", err)
	}

	return nil
}
