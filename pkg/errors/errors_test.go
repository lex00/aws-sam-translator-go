package errors

import (
	"errors"
	"testing"
)

func TestInvalidDocumentException_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "missing Resources section",
			expected: "invalid document: missing Resources section",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "invalid document: ",
		},
		{
			name:     "complex message",
			message:  "Transform must be 'AWS::Serverless-2016-10-31'",
			expected: "invalid document: Transform must be 'AWS::Serverless-2016-10-31'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidDocumentException{Message: tt.message}
			if got := err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInvalidDocumentException_ImplementsError(t *testing.T) {
	// Compile-time check that InvalidDocumentException implements error
	var _ error = (*InvalidDocumentException)(nil)

	err := &InvalidDocumentException{Message: "test"}
	if err.Error() == "" {
		t.Error("InvalidDocumentException.Error() should not return empty string")
	}
}

func TestInvalidResourceException_Error(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		message    string
		expected   string
	}{
		{
			name:       "simple resource error",
			resourceID: "MyFunction",
			message:    "Runtime is required",
			expected:   "invalid resource 'MyFunction': Runtime is required",
		},
		{
			name:       "empty resource ID",
			resourceID: "",
			message:    "missing type",
			expected:   "invalid resource '': missing type",
		},
		{
			name:       "complex resource name",
			resourceID: "MyApi/GET/users",
			message:    "invalid path",
			expected:   "invalid resource 'MyApi/GET/users': invalid path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidResourceException{
				ResourceID: tt.resourceID,
				Message:    tt.message,
			}
			if got := err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInvalidResourceException_ImplementsError(t *testing.T) {
	// Compile-time check that InvalidResourceException implements error
	var _ error = (*InvalidResourceException)(nil)

	err := &InvalidResourceException{ResourceID: "test", Message: "test"}
	if err.Error() == "" {
		t.Error("InvalidResourceException.Error() should not return empty string")
	}
}

func TestInvalidEventException_Error(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		eventID    string
		message    string
		expected   string
	}{
		{
			name:       "simple event error",
			resourceID: "MyFunction",
			eventID:    "ApiEvent",
			message:    "Path is required",
			expected:   "invalid event 'ApiEvent' on resource 'MyFunction': Path is required",
		},
		{
			name:       "empty event ID",
			resourceID: "MyFunction",
			eventID:    "",
			message:    "missing type",
			expected:   "invalid event '' on resource 'MyFunction': missing type",
		},
		{
			name:       "S3 event",
			resourceID: "ProcessorFunction",
			eventID:    "S3Upload",
			message:    "Bucket is required",
			expected:   "invalid event 'S3Upload' on resource 'ProcessorFunction': Bucket is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidEventException{
				ResourceID: tt.resourceID,
				EventID:    tt.eventID,
				Message:    tt.message,
			}
			if got := err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInvalidEventException_ImplementsError(t *testing.T) {
	// Compile-time check that InvalidEventException implements error
	var _ error = (*InvalidEventException)(nil)

	err := &InvalidEventException{ResourceID: "test", EventID: "test", Message: "test"}
	if err.Error() == "" {
		t.Error("InvalidEventException.Error() should not return empty string")
	}
}

func TestErrorsCanBeUsedWithErrorsIs(t *testing.T) {
	docErr := &InvalidDocumentException{Message: "test"}
	resErr := &InvalidResourceException{ResourceID: "res", Message: "test"}
	eventErr := &InvalidEventException{ResourceID: "res", EventID: "evt", Message: "test"}

	// These should all be usable as errors
	var _ error = docErr
	var _ error = resErr
	var _ error = eventErr

	// Verify they can be wrapped and unwrapped
	wrappedDoc := errors.Join(errors.New("context"), docErr)
	if wrappedDoc == nil {
		t.Error("should be able to wrap InvalidDocumentException")
	}

	wrappedRes := errors.Join(errors.New("context"), resErr)
	if wrappedRes == nil {
		t.Error("should be able to wrap InvalidResourceException")
	}

	wrappedEvent := errors.Join(errors.New("context"), eventErr)
	if wrappedEvent == nil {
		t.Error("should be able to wrap InvalidEventException")
	}
}
