// Package sam provides SAM resource transformers.
package sam

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// LayerVersion represents an AWS::Serverless::LayerVersion resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-layerversion.html
type LayerVersion struct {
	// ContentUri specifies the layer content location.
	// Can be a string (s3://bucket/key) or an object with Bucket, Key, Version properties.
	ContentUri interface{} `json:"ContentUri,omitempty" yaml:"ContentUri,omitempty"`

	// LayerName is the name of the layer. If not specified, the logical ID is used.
	LayerName string `json:"LayerName,omitempty" yaml:"LayerName,omitempty"`

	// Description is the description of the layer.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// CompatibleRuntimes is a list of compatible function runtimes.
	CompatibleRuntimes []string `json:"CompatibleRuntimes,omitempty" yaml:"CompatibleRuntimes,omitempty"`

	// CompatibleArchitectures is a list of compatible architectures.
	CompatibleArchitectures []string `json:"CompatibleArchitectures,omitempty" yaml:"CompatibleArchitectures,omitempty"`

	// LicenseInfo is the layer's software license information.
	LicenseInfo string `json:"LicenseInfo,omitempty" yaml:"LicenseInfo,omitempty"`

	// RetentionPolicy specifies whether to retain or delete the layer version.
	// Valid values: Retain, Delete (case-insensitive). Maps to DeletionPolicy.
	// If not specified, defaults to Retain.
	RetentionPolicy interface{} `json:"RetentionPolicy,omitempty" yaml:"RetentionPolicy,omitempty"`
}

// LayerVersionTransformer transforms AWS::Serverless::LayerVersion to CloudFormation.
type LayerVersionTransformer struct{}

// NewLayerVersionTransformer creates a new LayerVersionTransformer.
func NewLayerVersionTransformer() *LayerVersionTransformer {
	return &LayerVersionTransformer{}
}

// LayerVersionHashLength is the length of the hash suffix for layer logical IDs.
const LayerVersionHashLength = 10

// Transform converts a SAM LayerVersion to a CloudFormation Lambda::LayerVersion resource.
// It returns the resources map and the new logical ID (with hash suffix).
func (t *LayerVersionTransformer) Transform(logicalID string, lv *LayerVersion) (map[string]interface{}, string, error) {
	// Parse ContentUri to get S3 bucket and key
	content, err := t.parseContentUri(lv.ContentUri)
	if err != nil {
		return nil, "", err
	}

	// Build properties map
	properties := make(map[string]interface{})
	properties["Content"] = content

	// Use LayerName if specified, otherwise use the original logical ID
	if lv.LayerName != "" {
		properties["LayerName"] = lv.LayerName
	} else {
		properties["LayerName"] = logicalID
	}

	if lv.Description != "" {
		properties["Description"] = lv.Description
	}

	if len(lv.CompatibleRuntimes) > 0 {
		properties["CompatibleRuntimes"] = lv.CompatibleRuntimes
	}

	if len(lv.CompatibleArchitectures) > 0 {
		properties["CompatibleArchitectures"] = lv.CompatibleArchitectures
	}

	if lv.LicenseInfo != "" {
		properties["LicenseInfo"] = lv.LicenseInfo
	}

	// Determine DeletionPolicy from RetentionPolicy
	deletionPolicy := t.mapRetentionPolicy(lv.RetentionPolicy)

	// Generate new logical ID with hash suffix
	newLogicalID := t.generateLogicalID(logicalID, content)

	// Build the CloudFormation resource
	resource := map[string]interface{}{
		"Type":           "AWS::Lambda::LayerVersion",
		"DeletionPolicy": deletionPolicy,
		"Properties":     properties,
	}

	resources := map[string]interface{}{
		newLogicalID: resource,
	}

	return resources, newLogicalID, nil
}

// parseContentUri parses the ContentUri and returns the Content object for CloudFormation.
func (t *LayerVersionTransformer) parseContentUri(contentUri interface{}) (map[string]interface{}, error) {
	if contentUri == nil {
		return nil, fmt.Errorf("ContentUri is required for LayerVersion")
	}

	switch v := contentUri.(type) {
	case string:
		// Parse s3://bucket/key format
		return t.parseS3Uri(v)
	case map[string]interface{}:
		// Already an object with Bucket, Key, Version
		return t.parseContentUriObject(v)
	case map[interface{}]interface{}:
		// Convert to string keys (common in YAML parsing)
		converted := make(map[string]interface{})
		for key, val := range v {
			if strKey, ok := key.(string); ok {
				converted[strKey] = val
			}
		}
		return t.parseContentUriObject(converted)
	default:
		return nil, fmt.Errorf("invalid ContentUri type: %T", contentUri)
	}
}

// parseS3Uri parses an S3 URI string (s3://bucket/key) into Content properties.
func (t *LayerVersionTransformer) parseS3Uri(uri string) (map[string]interface{}, error) {
	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("invalid S3 URI: must start with s3://")
	}

	// Remove s3:// prefix
	path := strings.TrimPrefix(uri, "s3://")

	// Split into bucket and key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid S3 URI: must have bucket and key")
	}

	return map[string]interface{}{
		"S3Bucket": parts[0],
		"S3Key":    parts[1],
	}, nil
}

// parseContentUriObject parses a ContentUri object with Bucket, Key, Version properties.
func (t *LayerVersionTransformer) parseContentUriObject(obj map[string]interface{}) (map[string]interface{}, error) {
	content := make(map[string]interface{})

	if bucket, ok := obj["Bucket"]; ok {
		content["S3Bucket"] = bucket
	} else {
		return nil, fmt.Errorf("ContentUri object missing Bucket property")
	}

	if key, ok := obj["Key"]; ok {
		content["S3Key"] = key
	} else {
		return nil, fmt.Errorf("ContentUri object missing Key property")
	}

	if version, ok := obj["Version"]; ok {
		content["S3ObjectVersion"] = version
	}

	return content, nil
}

// mapRetentionPolicy maps SAM RetentionPolicy to CloudFormation DeletionPolicy.
// Default is Retain if not specified.
func (t *LayerVersionTransformer) mapRetentionPolicy(retentionPolicy interface{}) string {
	if retentionPolicy == nil {
		return "Retain"
	}

	// Handle string values
	if strPolicy, ok := retentionPolicy.(string); ok {
		// Case-insensitive comparison
		lowerPolicy := strings.ToLower(strPolicy)
		if lowerPolicy == "delete" {
			return "Delete"
		}
		return "Retain"
	}

	// If it's an intrinsic function (map), return Retain as default
	// The actual value will be resolved at deployment time
	return "Retain"
}

// generateLogicalID generates a new logical ID with a hash suffix based on content.
func (t *LayerVersionTransformer) generateLogicalID(baseID string, content map[string]interface{}) string {
	// Create hash input from content properties
	hashInput := fmt.Sprintf("%v%v", content["S3Bucket"], content["S3Key"])
	if version, ok := content["S3ObjectVersion"]; ok {
		hashInput += fmt.Sprintf("%v", version)
	}

	// Calculate hash
	h := sha256.New()
	h.Write([]byte(hashInput))
	hash := hex.EncodeToString(h.Sum(nil))[:LayerVersionHashLength]

	return baseID + hash
}
