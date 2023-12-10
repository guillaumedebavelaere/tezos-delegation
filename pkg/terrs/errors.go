package terrs

import (
	"fmt"
	"strings"
)

const message = "test error"

// TestError is used for error testing.
type TestError struct {
	messages []string
}

// NewTestError creates an TestError with messages.
func NewTestError(messages ...string) *TestError {
	return &TestError{messages: messages}
}

// Error is a function getting the error messages from the struct.
func (e *TestError) Error() string {
	if len(e.messages) == 0 {
		return message
	}

	return fmt.Sprintf("%s: %s", message, strings.Join(e.messages, ", "))
}
