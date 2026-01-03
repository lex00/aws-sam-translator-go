package translator

import "testing"

func TestNew(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Error("New() returned nil")
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version is empty")
	}
}
