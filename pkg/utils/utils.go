// Package utils provides utility functions for template processing.
package utils

import (
	"reflect"
	"sort"
)

// DeepMerge merges two maps recursively.
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range src {
		if dstV, ok := result[k]; ok {
			if dstMap, ok := dstV.(map[string]interface{}); ok {
				if srcMap, ok := v.(map[string]interface{}); ok {
					result[k] = DeepMerge(dstMap, srcMap)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}

// DeepCopy creates a deep copy of a map.
func DeepCopy(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			result[k] = DeepCopy(vm)
		} else if vs, ok := v.([]interface{}); ok {
			result[k] = deepCopySlice(vs)
		} else {
			result[k] = v
		}
	}
	return result
}

func deepCopySlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		if vm, ok := v.(map[string]interface{}); ok {
			result[i] = DeepCopy(vm)
		} else if vs, ok := v.([]interface{}); ok {
			result[i] = deepCopySlice(vs)
		} else {
			result[i] = v
		}
	}
	return result
}

// SortedKeys returns the keys of a map in sorted order.
// This is useful for deterministic iteration over maps.
func SortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedStringKeys returns the keys of a string-to-string map in sorted order.
func SortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedResourceKeys returns the keys of a resource map in sorted order.
// This provides deterministic output when iterating over resources.
func SortedResourceKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// DeepEqual compares two values for deep equality.
// It handles maps, slices, and primitive types.
func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// MapContains checks if all key-value pairs in subset exist in m.
func MapContains(m, subset map[string]interface{}) bool {
	for k, v := range subset {
		if mv, ok := m[k]; !ok || !DeepEqual(mv, v) {
			return false
		}
	}
	return true
}

// StringSliceContains checks if a string slice contains a specific value.
func StringSliceContains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}

// SortStringSlice returns a sorted copy of a string slice.
func SortStringSlice(slice []string) []string {
	result := make([]string, len(slice))
	copy(result, slice)
	sort.Strings(result)
	return result
}

// UniqueStrings returns a deduplicated, sorted slice of strings.
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	sort.Strings(result)
	return result
}

// MergeMaps merges multiple maps into a single map.
// Later maps take precedence over earlier ones.
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// GetStringValue safely gets a string value from a map.
// Returns the default value if the key doesn't exist or the value is not a string.
func GetStringValue(m map[string]interface{}, key, defaultValue string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

// GetMapValue safely gets a map value from a map.
// Returns nil if the key doesn't exist or the value is not a map.
func GetMapValue(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			return mv
		}
	}
	return nil
}

// GetSliceValue safely gets a slice value from a map.
// Returns nil if the key doesn't exist or the value is not a slice.
func GetSliceValue(m map[string]interface{}, key string) []interface{} {
	if v, ok := m[key]; ok {
		if sv, ok := v.([]interface{}); ok {
			return sv
		}
	}
	return nil
}
