package operation

import (
	"bytes"

	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type GoogleOperationError compute.OperationError
type GoogleOperationErrorB computebeta.OperationError

func (e GoogleOperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}

func (e GoogleOperationErrorB) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}
