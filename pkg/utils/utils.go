// Package utils provides utility functions for template processing.
package utils

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
