package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestResource_JSONSerialization(t *testing.T) {
	resource := Resource{
		ParentId:  "parent123",
		PathPart:  "users",
		RestApiId: "api123",
	}

	// Test JSON marshaling
	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource from JSON: %v", err)
	}

	if unmarshaled.ParentId != resource.ParentId {
		t.Errorf("ParentId mismatch: got %v, want %v", unmarshaled.ParentId, resource.ParentId)
	}

	if unmarshaled.PathPart != resource.PathPart {
		t.Errorf("PathPart mismatch: got %v, want %v", unmarshaled.PathPart, resource.PathPart)
	}

	if unmarshaled.RestApiId != resource.RestApiId {
		t.Errorf("RestApiId mismatch: got %v, want %v", unmarshaled.RestApiId, resource.RestApiId)
	}
}

func TestResource_YAMLSerialization(t *testing.T) {
	resource := Resource{
		ParentId:  "parent456",
		PathPart:  "{userId}",
		RestApiId: "api789",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Resource
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource from YAML: %v", err)
	}

	if unmarshaled.ParentId != resource.ParentId {
		t.Errorf("ParentId mismatch: got %v, want %v", unmarshaled.ParentId, resource.ParentId)
	}

	if unmarshaled.PathPart != resource.PathPart {
		t.Errorf("PathPart mismatch: got %v, want %v", unmarshaled.PathPart, resource.PathPart)
	}
}

func TestResource_WithIntrinsicFunctions(t *testing.T) {
	resource := Resource{
		ParentId:  map[string]interface{}{"Fn::GetAtt": []string{"RestApi", "RootResourceId"}},
		PathPart:  "users",
		RestApiId: map[string]interface{}{"Ref": "RestApi"},
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource with intrinsic functions: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	parentIdMap, ok := unmarshaled.ParentId.(map[string]interface{})
	if !ok {
		t.Error("ParentId should be a map for intrinsic function")
	} else if _, exists := parentIdMap["Fn::GetAtt"]; !exists {
		t.Error("ParentId should contain Fn::GetAtt intrinsic function")
	}

	restApiIdMap, ok := unmarshaled.RestApiId.(map[string]interface{})
	if !ok {
		t.Error("RestApiId should be a map for intrinsic function")
	} else if _, exists := restApiIdMap["Ref"]; !exists {
		t.Error("RestApiId should contain Ref intrinsic function")
	}
}

func TestResource_WithPathParameter(t *testing.T) {
	resource := Resource{
		ParentId:  "parent123",
		PathPart:  "{id}",
		RestApiId: "api123",
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource with path parameter: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource with path parameter: %v", err)
	}

	if unmarshaled.PathPart != "{id}" {
		t.Errorf("PathPart mismatch: got %v, want {id}", unmarshaled.PathPart)
	}
}

func TestResource_ProxyResource(t *testing.T) {
	resource := Resource{
		ParentId:  "parent123",
		PathPart:  "{proxy+}",
		RestApiId: "api123",
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal proxy Resource: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal proxy Resource: %v", err)
	}

	if unmarshaled.PathPart != "{proxy+}" {
		t.Errorf("PathPart mismatch: got %v, want {proxy+}", unmarshaled.PathPart)
	}
}

func TestResource_NestedPath(t *testing.T) {
	// Test a nested resource structure
	parentResource := Resource{
		ParentId:  "root123",
		PathPart:  "api",
		RestApiId: "api123",
	}

	childResource := Resource{
		ParentId:  "parent123",
		PathPart:  "v1",
		RestApiId: "api123",
	}

	// Marshal parent
	parentData, err := json.Marshal(parentResource)
	if err != nil {
		t.Fatalf("Failed to marshal parent Resource: %v", err)
	}

	// Marshal child
	childData, err := json.Marshal(childResource)
	if err != nil {
		t.Fatalf("Failed to marshal child Resource: %v", err)
	}

	// Unmarshal and verify
	var unmarshaledParent Resource
	if err := json.Unmarshal(parentData, &unmarshaledParent); err != nil {
		t.Fatalf("Failed to unmarshal parent Resource: %v", err)
	}

	var unmarshaledChild Resource
	if err := json.Unmarshal(childData, &unmarshaledChild); err != nil {
		t.Fatalf("Failed to unmarshal child Resource: %v", err)
	}

	if unmarshaledParent.PathPart != "api" {
		t.Errorf("Parent PathPart mismatch: got %v, want api", unmarshaledParent.PathPart)
	}

	if unmarshaledChild.PathPart != "v1" {
		t.Errorf("Child PathPart mismatch: got %v, want v1", unmarshaledChild.PathPart)
	}
}
