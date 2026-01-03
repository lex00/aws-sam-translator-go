package lambda

// LayerVersion represents an AWS::Lambda::LayerVersion CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html
type LayerVersion struct {
	// CompatibleArchitectures is a list of compatible architectures.
	// Valid values: x86_64, arm64
	CompatibleArchitectures []string `json:"CompatibleArchitectures,omitempty" yaml:"CompatibleArchitectures,omitempty"`

	// CompatibleRuntimes is a list of compatible function runtimes.
	CompatibleRuntimes []string `json:"CompatibleRuntimes,omitempty" yaml:"CompatibleRuntimes,omitempty"`

	// Content is the layer content (required).
	Content *LayerContent `json:"Content" yaml:"Content"`

	// Description is the layer description.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// LayerName is the name of the layer.
	LayerName string `json:"LayerName,omitempty" yaml:"LayerName,omitempty"`

	// LicenseInfo is the layer's software license.
	LicenseInfo string `json:"LicenseInfo,omitempty" yaml:"LicenseInfo,omitempty"`
}

// LayerContent represents the content of a Lambda layer.
type LayerContent struct {
	// S3Bucket is the S3 bucket name containing the layer content.
	S3Bucket interface{} `json:"S3Bucket,omitempty" yaml:"S3Bucket,omitempty"`

	// S3Key is the S3 object key of the layer content.
	S3Key interface{} `json:"S3Key,omitempty" yaml:"S3Key,omitempty"`

	// S3ObjectVersion is the S3 object version of the layer content.
	S3ObjectVersion interface{} `json:"S3ObjectVersion,omitempty" yaml:"S3ObjectVersion,omitempty"`
}

// NewLayerVersion creates a new LayerVersion with the required content.
func NewLayerVersion(content *LayerContent) *LayerVersion {
	return &LayerVersion{
		Content: content,
	}
}

// NewLayerVersionFromS3 creates a LayerVersion from an S3 location.
func NewLayerVersionFromS3(bucket, key interface{}) *LayerVersion {
	return &LayerVersion{
		Content: &LayerContent{
			S3Bucket: bucket,
			S3Key:    key,
		},
	}
}

// NewLayerVersionFromS3WithVersion creates a LayerVersion from an S3 location with version.
func NewLayerVersionFromS3WithVersion(bucket, key, version interface{}) *LayerVersion {
	return &LayerVersion{
		Content: &LayerContent{
			S3Bucket:        bucket,
			S3Key:           key,
			S3ObjectVersion: version,
		},
	}
}

// WithLayerName sets the layer name.
func (l *LayerVersion) WithLayerName(name string) *LayerVersion {
	l.LayerName = name
	return l
}

// WithDescription sets the layer description.
func (l *LayerVersion) WithDescription(description string) *LayerVersion {
	l.Description = description
	return l
}

// WithLicenseInfo sets the layer's license information.
func (l *LayerVersion) WithLicenseInfo(license string) *LayerVersion {
	l.LicenseInfo = license
	return l
}

// WithCompatibleRuntimes sets the compatible runtimes.
func (l *LayerVersion) WithCompatibleRuntimes(runtimes ...string) *LayerVersion {
	l.CompatibleRuntimes = runtimes
	return l
}

// WithCompatibleArchitectures sets the compatible architectures.
func (l *LayerVersion) WithCompatibleArchitectures(architectures ...string) *LayerVersion {
	l.CompatibleArchitectures = architectures
	return l
}

// AddCompatibleRuntime adds a compatible runtime.
func (l *LayerVersion) AddCompatibleRuntime(runtime string) *LayerVersion {
	l.CompatibleRuntimes = append(l.CompatibleRuntimes, runtime)
	return l
}

// AddCompatibleArchitecture adds a compatible architecture.
func (l *LayerVersion) AddCompatibleArchitecture(architecture string) *LayerVersion {
	l.CompatibleArchitectures = append(l.CompatibleArchitectures, architecture)
	return l
}

// ToCloudFormation converts the LayerVersion to a CloudFormation resource.
func (l *LayerVersion) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	if l.Content != nil {
		properties["Content"] = l.Content.toMap()
	}

	if len(l.CompatibleArchitectures) > 0 {
		properties["CompatibleArchitectures"] = l.CompatibleArchitectures
	}
	if len(l.CompatibleRuntimes) > 0 {
		properties["CompatibleRuntimes"] = l.CompatibleRuntimes
	}
	if l.Description != "" {
		properties["Description"] = l.Description
	}
	if l.LayerName != "" {
		properties["LayerName"] = l.LayerName
	}
	if l.LicenseInfo != "" {
		properties["LicenseInfo"] = l.LicenseInfo
	}

	return map[string]interface{}{
		"Type":       ResourceTypeLayerVersion,
		"Properties": properties,
	}
}

func (c *LayerContent) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if c.S3Bucket != nil {
		m["S3Bucket"] = c.S3Bucket
	}
	if c.S3Key != nil {
		m["S3Key"] = c.S3Key
	}
	if c.S3ObjectVersion != nil {
		m["S3ObjectVersion"] = c.S3ObjectVersion
	}
	return m
}

// LayerVersionPermission represents an AWS::Lambda::LayerVersionPermission CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html
type LayerVersionPermission struct {
	// Action is the API action that grants access (required).
	// Typically "lambda:GetLayerVersion".
	Action string `json:"Action" yaml:"Action"`

	// LayerVersionArn is the ARN of the layer version (required).
	LayerVersionArn interface{} `json:"LayerVersionArn" yaml:"LayerVersionArn"`

	// OrganizationId restricts access to accounts in the organization.
	OrganizationId string `json:"OrganizationId,omitempty" yaml:"OrganizationId,omitempty"`

	// Principal is the account ID, account pattern, or AWS service that can access (required).
	Principal string `json:"Principal" yaml:"Principal"`
}

// ResourceTypeLayerVersionPermission is the CloudFormation resource type.
const ResourceTypeLayerVersionPermission = "AWS::Lambda::LayerVersionPermission"

// NewLayerVersionPermission creates a new LayerVersionPermission.
func NewLayerVersionPermission(layerVersionArn interface{}, principal string) *LayerVersionPermission {
	return &LayerVersionPermission{
		Action:          "lambda:GetLayerVersion",
		LayerVersionArn: layerVersionArn,
		Principal:       principal,
	}
}

// NewLayerVersionPermissionPublic creates a public LayerVersionPermission.
func NewLayerVersionPermissionPublic(layerVersionArn interface{}) *LayerVersionPermission {
	return &LayerVersionPermission{
		Action:          "lambda:GetLayerVersion",
		LayerVersionArn: layerVersionArn,
		Principal:       "*",
	}
}

// NewLayerVersionPermissionOrg creates a LayerVersionPermission for an organization.
func NewLayerVersionPermissionOrg(layerVersionArn interface{}, organizationId string) *LayerVersionPermission {
	return &LayerVersionPermission{
		Action:          "lambda:GetLayerVersion",
		LayerVersionArn: layerVersionArn,
		Principal:       "*",
		OrganizationId:  organizationId,
	}
}

// WithOrganizationId sets the organization ID restriction.
func (p *LayerVersionPermission) WithOrganizationId(orgId string) *LayerVersionPermission {
	p.OrganizationId = orgId
	return p
}

// ToCloudFormation converts the LayerVersionPermission to a CloudFormation resource.
func (p *LayerVersionPermission) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["Action"] = p.Action
	properties["LayerVersionArn"] = p.LayerVersionArn
	properties["Principal"] = p.Principal

	if p.OrganizationId != "" {
		properties["OrganizationId"] = p.OrganizationId
	}

	return map[string]interface{}{
		"Type":       ResourceTypeLayerVersionPermission,
		"Properties": properties,
	}
}
