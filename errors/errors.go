package errors

import (
	"fmt"
)

type Error interface {
	error
}

type UnwrappableError interface {
	Error
	Unwrap() error
}

type BuildTransaction struct {
	Err error
}

var _ UnwrappableError = (*BuildTransaction)(nil)

func NewBuildTransaction(err error) BuildTransaction {
	return BuildTransaction{Err: err}
}

func (b BuildTransaction) Error() string {
	return fmt.Sprintf("error building transaction: %v", b.Err)
}

func (b BuildTransaction) Unwrap() error {
	return b.Err
}
