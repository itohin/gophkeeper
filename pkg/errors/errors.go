package errors

import "fmt"

type InvalidArgumentError struct {
	Err error
}

func NewInvalidArgumentError(err error) error {
	return &InvalidArgumentError{Err: err}
}

func (i *InvalidArgumentError) Error() string {
	return fmt.Sprintf("%v", i.Err)
}

type DomainError struct {
	Err error
}

func NewDomainError(err error) error {
	return &DomainError{Err: err}
}

func (i *DomainError) Error() string {
	return fmt.Sprintf("%v", i.Err)
}

type AuthError struct {
	Err error
}

func NewAuthError(err error) error {
	return &AuthError{Err: err}
}

func (i *AuthError) Error() string {
	return fmt.Sprintf("%v", i.Err)
}
