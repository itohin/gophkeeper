package usecases

import "fmt"

type InvalidArgumentError struct {
	Err error
}

func NewInvalidArgumentError(err error) error {
	return &InvalidArgumentError{Err: err}
}

func (i *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument error: %v", i.Err)
}
