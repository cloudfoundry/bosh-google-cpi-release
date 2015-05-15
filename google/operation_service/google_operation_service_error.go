package operation

import (
	"bytes"

	"google.golang.org/api/compute/v1"
)

type GoogleOperationError compute.OperationError

func (e GoogleOperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}
