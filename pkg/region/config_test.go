package region

import (
	"reflect"
	"testing"
)

func TestGetPartitionForRegion(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		expected Partition
	}{
		{
			name:     "empty region defaults to AWS",
			region:   "",
			expected: PartitionAWS,
		},
		{
			name:     "us-east-1",
			region:   "us-east-1",
			expected: PartitionAWS,
		},
		{
			name:     "us-west-2",
			region:   "us-west-2",
			expected: PartitionAWS,
		},
		{
			name:     "eu-west-1",
			region:   "eu-west-1",
			expected: PartitionAWS,
		},
		{
			name:     "ap-southeast-1",
			region:   "ap-southeast-1",
			expected: PartitionAWS,
		},
		{
			name:     "sa-east-1",
			region:   "sa-east-1",
			expected: PartitionAWS,
		},
		{
			name:     "cn-north-1 (China)",
			region:   "cn-north-1",
			expected: PartitionAWSChina,
		},
		{
			name:     "cn-northwest-1 (China)",
			region:   "cn-northwest-1",
			expected: PartitionAWSChina,
		},
		{
			name:     "us-gov-west-1 (GovCloud)",
			region:   "us-gov-west-1",
			expected: PartitionAWSGov,
		},
		{
			name:     "us-gov-east-1 (GovCloud)",
			region:   "us-gov-east-1",
			expected: PartitionAWSGov,
		},
		{
			name:     "unknown region defaults to AWS",
			region:   "unknown-region-1",
			expected: PartitionAWS,
		},
		{
			name:     "me-south-1 (Middle East)",
			region:   "me-south-1",
			expected: PartitionAWS,
		},
		{
			name:     "af-south-1 (Africa)",
			region:   "af-south-1",
			expected: PartitionAWS,
		},
		{
			name:     "il-central-1 (Israel)",
			region:   "il-central-1",
			expected: PartitionAWS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPartitionForRegion(tt.region)
			if result != tt.expected {
				t.Errorf("GetPartitionForRegion(%q) = %v, want %v", tt.region, result, tt.expected)
			}
		})
	}
}

func TestGetPartitionConfig(t *testing.T) {
	tests := []struct {
		name              string
		partition         Partition
		expectedDNSSuffix string
	}{
		{
			name:              "AWS partition",
			partition:         PartitionAWS,
			expectedDNSSuffix: "amazonaws.com",
		},
		{
			name:              "AWS China partition",
			partition:         PartitionAWSChina,
			expectedDNSSuffix: "amazonaws.com.cn",
		},
		{
			name:              "AWS GovCloud partition",
			partition:         PartitionAWSGov,
			expectedDNSSuffix: "amazonaws.com",
		},
		{
			name:              "unknown partition defaults to AWS",
			partition:         Partition("unknown"),
			expectedDNSSuffix: "amazonaws.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetPartitionConfig(tt.partition)
			if config.DNSSuffix != tt.expectedDNSSuffix {
				t.Errorf("GetPartitionConfig(%v).DNSSuffix = %v, want %v",
					tt.partition, config.DNSSuffix, tt.expectedDNSSuffix)
			}
		})
	}
}

func TestGetDNSSuffix(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		expected string
	}{
		{
			name:     "us-east-1",
			region:   "us-east-1",
			expected: "amazonaws.com",
		},
		{
			name:     "cn-north-1",
			region:   "cn-north-1",
			expected: "amazonaws.com.cn",
		},
		{
			name:     "us-gov-west-1",
			region:   "us-gov-west-1",
			expected: "amazonaws.com",
		},
		{
			name:     "empty region",
			region:   "",
			expected: "amazonaws.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDNSSuffix(tt.region)
			if result != tt.expected {
				t.Errorf("GetDNSSuffix(%q) = %v, want %v", tt.region, result, tt.expected)
			}
		})
	}
}

func TestGetArnPartition(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		expected string
	}{
		{
			name:     "us-east-1",
			region:   "us-east-1",
			expected: "aws",
		},
		{
			name:     "cn-north-1",
			region:   "cn-north-1",
			expected: "aws-cn",
		},
		{
			name:     "us-gov-west-1",
			region:   "us-gov-west-1",
			expected: "aws-us-gov",
		},
		{
			name:     "empty region",
			region:   "",
			expected: "aws",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetArnPartition(tt.region)
			if result != tt.expected {
				t.Errorf("GetArnPartition(%q) = %v, want %v", tt.region, result, tt.expected)
			}
		})
	}
}

func TestIsValidRegion(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		expected bool
	}{
		{
			name:     "valid us-east-1",
			region:   "us-east-1",
			expected: true,
		},
		{
			name:     "valid eu-west-1",
			region:   "eu-west-1",
			expected: true,
		},
		{
			name:     "valid ap-southeast-1",
			region:   "ap-southeast-1",
			expected: true,
		},
		{
			name:     "valid cn-north-1",
			region:   "cn-north-1",
			expected: true,
		},
		{
			name:     "valid us-gov-west-1",
			region:   "us-gov-west-1",
			expected: true,
		},
		{
			name:     "empty region",
			region:   "",
			expected: false,
		},
		{
			name:     "no hyphens",
			region:   "useast1",
			expected: false,
		},
		{
			name:     "single segment",
			region:   "us",
			expected: false,
		},
		{
			name:     "two segments only",
			region:   "us-east",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidRegion(tt.region)
			if result != tt.expected {
				t.Errorf("IsValidRegion(%q) = %v, want %v", tt.region, result, tt.expected)
			}
		})
	}
}

func TestRegionOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		expected string
	}{
		{
			name:     "empty region returns default",
			region:   "",
			expected: DefaultRegion,
		},
		{
			name:     "specified region returned",
			region:   "eu-west-1",
			expected: "eu-west-1",
		},
		{
			name:     "us-east-1 returned as-is",
			region:   "us-east-1",
			expected: "us-east-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RegionOrDefault(tt.region)
			if result != tt.expected {
				t.Errorf("RegionOrDefault(%q) = %v, want %v", tt.region, result, tt.expected)
			}
		})
	}
}

func TestDefaultRegion(t *testing.T) {
	if DefaultRegion != "us-east-1" {
		t.Errorf("DefaultRegion = %v, want us-east-1", DefaultRegion)
	}
}

func TestAllPartitions(t *testing.T) {
	partitions := AllPartitions()

	// Check that all known partitions are returned
	expected := []Partition{PartitionAWS, PartitionAWSChina, PartitionAWSGov}

	if !reflect.DeepEqual(partitions, expected) {
		t.Errorf("AllPartitions() = %v, want %v", partitions, expected)
	}

	// Check partition constants
	if PartitionAWS != "aws" {
		t.Errorf("PartitionAWS = %v, want aws", PartitionAWS)
	}
	if PartitionAWSChina != "aws-cn" {
		t.Errorf("PartitionAWSChina = %v, want aws-cn", PartitionAWSChina)
	}
	if PartitionAWSGov != "aws-us-gov" {
		t.Errorf("PartitionAWSGov = %v, want aws-us-gov", PartitionAWSGov)
	}
}

func TestPartitionConfigRegexPatterns(t *testing.T) {
	// Ensure partition configs have regex patterns defined
	for _, partition := range AllPartitions() {
		config := GetPartitionConfig(partition)
		if config.RegionRegex == "" {
			t.Errorf("Partition %v has empty RegionRegex", partition)
		}
		if config.Partition != partition {
			t.Errorf("Partition config mismatch: got %v, want %v", config.Partition, partition)
		}
	}
}
