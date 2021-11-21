package errors

import "fmt"

type NotFoundError struct {
	Id string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("process %s not found", e.Id)
}
