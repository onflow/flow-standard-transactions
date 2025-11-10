package errors

import (
	"fmt"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
)

type Error interface {
	error
}

type UnwrappableError interface {
	Error
	Unwrap() error
}

type PrepareTransactionFromDNA struct {
	DNA models.DNA
	Err error
}

var _ UnwrappableError = (*PrepareTransactionFromDNA)(nil)

func NewPrepareTransactionFromDNA(err error, dna models.DNA) PrepareTransactionFromDNA {
	return PrepareTransactionFromDNA{Err: err, DNA: dna}
}

func (g PrepareTransactionFromDNA) Error() string {
	return fmt.Sprintf("error generating transaction from DNA  %s: %v", g.DNA, g.Err)
}

func (g PrepareTransactionFromDNA) Unwrap() error {
	return g.Err
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
