// Package region provides AWS region and partition configuration handling.
package region

import (
	"strings"
)

// Partition represents an AWS partition.
type Partition string

const (
	// PartitionAWS is the standard AWS partition.
	PartitionAWS Partition = "aws"
	// PartitionAWSChina is the AWS China partition.
	PartitionAWSChina Partition = "aws-cn"
	// PartitionAWSGov is the AWS GovCloud partition.
	PartitionAWSGov Partition = "aws-us-gov"
)

// PartitionConfig contains configuration for an AWS partition.
type PartitionConfig struct {
	// Partition is the partition identifier.
	Partition Partition
	// DNSSuffix is the DNS suffix for the partition (e.g., amazonaws.com).
	DNSSuffix string
	// RegionRegex is a regex pattern that matches regions in this partition.
	RegionRegex string
}

// partitionConfigs maps partition identifiers to their configuration.
var partitionConfigs = map[Partition]PartitionConfig{
	PartitionAWS: {
		Partition:   PartitionAWS,
		DNSSuffix:   "amazonaws.com",
		RegionRegex: `^(us|eu|ap|sa|ca|me|af|il)\-\w+\-\d+$`,
	},
	PartitionAWSChina: {
		Partition:   PartitionAWSChina,
		DNSSuffix:   "amazonaws.com.cn",
		RegionRegex: `^cn\-\w+\-\d+$`,
	},
	PartitionAWSGov: {
		Partition:   PartitionAWSGov,
		DNSSuffix:   "amazonaws.com",
		RegionRegex: `^us\-gov\-\w+\-\d+$`,
	},
}

// regionPrefixToPartition maps region prefixes to partitions.
var regionPrefixToPartition = map[string]Partition{
	"cn-":     PartitionAWSChina,
	"us-gov-": PartitionAWSGov,
}

// DefaultRegion is the default AWS region used when none is specified.
const DefaultRegion = "us-east-1"

// GetPartitionForRegion returns the partition for a given region.
// Returns PartitionAWS if the region is empty or doesn't match any known partition.
func GetPartitionForRegion(region string) Partition {
	if region == "" {
		return PartitionAWS
	}

	// Check for specific partition prefixes
	for prefix, partition := range regionPrefixToPartition {
		if strings.HasPrefix(region, prefix) {
			return partition
		}
	}

	// Default to AWS partition
	return PartitionAWS
}

// GetPartitionConfig returns the configuration for a partition.
func GetPartitionConfig(partition Partition) PartitionConfig {
	if config, ok := partitionConfigs[partition]; ok {
		return config
	}
	// Default to AWS partition config
	return partitionConfigs[PartitionAWS]
}

// GetDNSSuffix returns the DNS suffix for a region.
func GetDNSSuffix(region string) string {
	partition := GetPartitionForRegion(region)
	config := GetPartitionConfig(partition)
	return config.DNSSuffix
}

// GetArnPartition returns the ARN partition string for a region.
// This is the partition identifier used in ARNs (e.g., "aws", "aws-cn", "aws-us-gov").
func GetArnPartition(region string) string {
	partition := GetPartitionForRegion(region)
	return string(partition)
}

// IsValidRegion checks if a region string matches the expected format.
// It does not validate that the region actually exists, only the format.
func IsValidRegion(region string) bool {
	if region == "" {
		return false
	}

	// Basic format check: regions should have at least one hyphen
	if !strings.Contains(region, "-") {
		return false
	}

	// Region should have format like: xx-xxxx-N or xxx-xxxx-N
	parts := strings.Split(region, "-")
	return len(parts) >= 3
}

// RegionOrDefault returns the region if it's not empty, otherwise returns the default region.
func RegionOrDefault(region string) string {
	if region == "" {
		return DefaultRegion
	}
	return region
}

// AllPartitions returns all known partitions.
func AllPartitions() []Partition {
	return []Partition{PartitionAWS, PartitionAWSChina, PartitionAWSGov}
}
