// Package errors provides error types for SAM template processing.
package errors

import "fmt"

// InvalidDocumentException is returned when a SAM template is malformed.
type InvalidDocumentException struct {
	Message string
}

func (e *InvalidDocumentException) Error() string {
	return fmt.Sprintf("invalid document: %s", e.Message)
}

// InvalidResourceException is returned when a SAM resource is invalid.
type InvalidResourceException struct {
	ResourceID string
	Message    string
}

func (e *InvalidResourceException) Error() string {
	return fmt.Sprintf("invalid resource '%s': %s", e.ResourceID, e.Message)
}

// InvalidEventException is returned when an event source configuration is invalid.
type InvalidEventException struct {
	ResourceID string
	EventID    string
	Message    string
}

func (e *InvalidEventException) Error() string {
	return fmt.Sprintf("invalid event '%s' on resource '%s': %s", e.EventID, e.ResourceID, e.Message)
}
