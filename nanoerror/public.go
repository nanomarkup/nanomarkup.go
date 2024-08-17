package nanoerror

import (
	"fmt"
	"strings"
)

// InvalidArgumentError describes an error that occurs when an invalid argument is provided.
type InvalidArgumentError struct {
	Context string
	Err     error
}

// InvalidEntityError describes an error that occurs when an attempt is made with an invalid entity.
type InvalidEntityError struct {
	Context string
	Entity  string
	Err     error
}

const (
	ErrorFmt string = "[%s] %s"
)

// Error returns a string representation of the InvalidArgumentError.
func (e *InvalidArgumentError) Error() string {
	if len(strings.TrimSpace(e.Context)) > 0 {
		return fmt.Sprintf(ErrorFmt, e.Context, e.Err.Error())
	} else {
		return e.Err.Error()
	}
}

// Error returns a string representation of the InvalidEntityError.
func (e *InvalidEntityError) Error() string {
	var s string
	if len(strings.TrimSpace(e.Entity)) > 0 {
		s = fmt.Sprintf("%s: %s", e.Err.Error(), e.Entity)
	} else {
		s = e.Err.Error()
	}

	if len(strings.TrimSpace(e.Context)) > 0 {
		return fmt.Sprintf(ErrorFmt, e.Context, s)
	} else {
		return s
	}
}
