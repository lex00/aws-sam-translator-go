package intrinsics

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewDependencyTracker(t *testing.T) {
	tracker := NewDependencyTracker()

	if tracker == nil {
		t.Fatal("expected tracker to be non-nil")
	}
	if tracker.dependencies == nil {
		t.Error("dependencies map should be initialized")
	}
	if tracker.dependents == nil {
		t.Error("dependents map should be initialized")
	}
}

func TestDependencyTrackerAddDependency(t *testing.T) {
	tracker := NewDependencyTracker()

	tracker.AddDependency("FunctionA", "BucketA")
	tracker.AddDependency("FunctionA", "TableA")
	tracker.AddDependency("FunctionB", "BucketA")

	// Check FunctionA dependencies
	depsA := tracker.GetDependencies("FunctionA")
	if len(depsA) != 2 {
		t.Errorf("expected 2 dependencies for FunctionA, got %d", len(depsA))
	}
	sort.Strings(depsA)
	if depsA[0] != "BucketA" || depsA[1] != "TableA" {
		t.Errorf("unexpected dependencies: %v", depsA)
	}

	// Check BucketA dependents
	dependentsB := tracker.GetDependents("BucketA")
	if len(dependentsB) != 2 {
		t.Errorf("expected 2 dependents for BucketA, got %d", len(dependentsB))
	}
}

func TestDependencyTrackerHasDependency(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("A", "B")

	if !tracker.HasDependency("A", "B") {
		t.Error("expected A to depend on B")
	}
	if tracker.HasDependency("B", "A") {
		t.Error("expected B NOT to depend on A")
	}
	if tracker.HasDependency("C", "D") {
		t.Error("expected non-existent dependency to return false")
	}
}

func TestDependencyTrackerAllDependencies(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("A", "B")
	tracker.AddDependency("A", "C")
	tracker.AddDependency("D", "E")

	all := tracker.AllDependencies()

	if len(all) != 2 {
		t.Errorf("expected 2 resources with dependencies, got %d", len(all))
	}

	depsA := all["A"]
	sort.Strings(depsA)
	expected := []string{"B", "C"}
	if !reflect.DeepEqual(depsA, expected) {
		t.Errorf("expected A deps %v, got %v", expected, depsA)
	}
}

func TestDependencyTrackerTopologicalSort(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("C", "B")
	tracker.AddDependency("B", "A")

	resources := []string{"A", "B", "C"}
	sorted, err := tracker.TopologicalSort(resources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// A should come before B, B before C
	indexA := indexOf(sorted, "A")
	indexB := indexOf(sorted, "B")
	indexC := indexOf(sorted, "C")

	if indexA > indexB || indexB > indexC {
		t.Errorf("expected order A -> B -> C, got %v", sorted)
	}
}

func TestDependencyTrackerTopologicalSortCircularDependency(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("A", "B")
	tracker.AddDependency("B", "C")
	tracker.AddDependency("C", "A") // Creates cycle

	resources := []string{"A", "B", "C"}
	_, err := tracker.TopologicalSort(resources)
	if err == nil {
		t.Fatal("expected circular dependency error")
	}

	_, ok := err.(*CircularDependencyError)
	if !ok {
		t.Errorf("expected CircularDependencyError, got %T", err)
	}
}

func TestDependencyTrackerMergeDependencies(t *testing.T) {
	tracker1 := NewDependencyTracker()
	tracker1.AddDependency("A", "B")

	tracker2 := NewDependencyTracker()
	tracker2.AddDependency("C", "D")
	tracker2.AddDependency("A", "E")

	tracker1.MergeDependencies(tracker2)

	// Check merged dependencies
	if !tracker1.HasDependency("A", "B") {
		t.Error("original dependency should remain")
	}
	if !tracker1.HasDependency("C", "D") {
		t.Error("merged dependency C->D should exist")
	}
	if !tracker1.HasDependency("A", "E") {
		t.Error("merged dependency A->E should exist")
	}
}

func TestDependencyTrackerClear(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("A", "B")
	tracker.AddDependency("C", "D")

	tracker.Clear()

	if tracker.Count() != 0 {
		t.Error("tracker should be empty after clear")
	}
}

func TestDependencyTrackerCount(t *testing.T) {
	tracker := NewDependencyTracker()

	if tracker.Count() != 0 {
		t.Error("empty tracker should have count 0")
	}

	tracker.AddDependency("A", "B")
	tracker.AddDependency("A", "C")
	tracker.AddDependency("D", "E")

	if tracker.Count() != 3 {
		t.Errorf("expected count 3, got %d", tracker.Count())
	}
}

func TestCircularDependencyError(t *testing.T) {
	err := &CircularDependencyError{Resource: "MyResource"}
	expected := "circular dependency detected involving resource: MyResource"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestResourceRefCollector(t *testing.T) {
	collector := NewResourceRefCollector()

	value := map[string]interface{}{
		"Properties": map[string]interface{}{
			"BucketName": map[string]interface{}{"Ref": "MyBucket"},
			"TableArn": map[string]interface{}{
				"Fn::GetAtt": []interface{}{"MyTable", "Arn"},
			},
			"Nested": map[string]interface{}{
				"Inner": map[string]interface{}{"Ref": "AnotherResource"},
			},
		},
	}

	refs := collector.CollectRefs(value, "Root")

	if len(refs) != 3 {
		t.Fatalf("expected 3 refs, got %d", len(refs))
	}

	// Find specific refs
	foundBucket := false
	foundTable := false
	foundAnother := false

	for _, ref := range refs {
		switch ref.LogicalID {
		case "MyBucket":
			foundBucket = true
			if ref.Attribute != "" {
				t.Error("Ref should not have an attribute")
			}
		case "MyTable":
			foundTable = true
			if ref.Attribute != "Arn" {
				t.Errorf("expected attribute 'Arn', got %q", ref.Attribute)
			}
		case "AnotherResource":
			foundAnother = true
		}
	}

	if !foundBucket {
		t.Error("MyBucket ref not found")
	}
	if !foundTable {
		t.Error("MyTable ref not found")
	}
	if !foundAnother {
		t.Error("AnotherResource ref not found")
	}
}

func TestResourceRefCollectorGetAttShortForm(t *testing.T) {
	collector := NewResourceRefCollector()

	value := map[string]interface{}{
		"Fn::GetAtt": "MyResource.MyAttribute",
	}

	refs := collector.CollectRefs(value, "Root")

	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}

	if refs[0].LogicalID != "MyResource" {
		t.Errorf("expected LogicalID 'MyResource', got %q", refs[0].LogicalID)
	}
	if refs[0].Attribute != "MyAttribute" {
		t.Errorf("expected Attribute 'MyAttribute', got %q", refs[0].Attribute)
	}
}

func TestResourceRefCollectorSlice(t *testing.T) {
	collector := NewResourceRefCollector()

	value := []interface{}{
		map[string]interface{}{"Ref": "Resource1"},
		map[string]interface{}{"Ref": "Resource2"},
	}

	refs := collector.CollectRefs(value, "Root")

	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestSplitFirst(t *testing.T) {
	tests := []struct {
		input    string
		sep      string
		expected []string
	}{
		{"a.b.c", ".", []string{"a", "b.c"}},
		{"no-separator", ".", []string{"no-separator"}},
		{"Resource.Attribute", ".", []string{"Resource", "Attribute"}},
	}

	for _, tt := range tests {
		result := splitFirst(tt.input, tt.sep)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("splitFirst(%q, %q) = %v, expected %v", tt.input, tt.sep, result, tt.expected)
		}
	}
}

func TestDependencyTrackerMergeNil(t *testing.T) {
	tracker := NewDependencyTracker()
	tracker.AddDependency("A", "B")

	// Merging nil should not panic or change anything
	tracker.MergeDependencies(nil)

	if !tracker.HasDependency("A", "B") {
		t.Error("dependency should still exist after merging nil")
	}
}

func TestGetDependenciesEmpty(t *testing.T) {
	tracker := NewDependencyTracker()

	deps := tracker.GetDependencies("NonExistent")
	if deps != nil {
		t.Errorf("expected nil for non-existent resource, got %v", deps)
	}
}

func TestGetDependentsEmpty(t *testing.T) {
	tracker := NewDependencyTracker()

	deps := tracker.GetDependents("NonExistent")
	if deps != nil {
		t.Errorf("expected nil for non-existent resource, got %v", deps)
	}
}

// Helper function to find index in slice
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
