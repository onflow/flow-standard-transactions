package errors

import "fmt"

type CardinalityMismatch struct {
	Expected uint
	Actual   uint
}

var _ Error = (*CardinalityMismatch)(nil)

func NewCardinalityMismatch(expected, actual uint) CardinalityMismatch {
	return CardinalityMismatch{Expected: expected, Actual: actual}
}

func (i CardinalityMismatch) Error() string {
	return fmt.Sprintf("cardinality mismatch, expected %d, got %d", i.Expected, i.Actual)
}

type NotImplemented struct {
	Feature string
}

var _ Error = (*NotImplemented)(nil)

func NewNotImplemented(feature string) NotImplemented {
	return NotImplemented{Feature: feature}
}

func (n NotImplemented) Error() string {
	return fmt.Sprintf("feature %s not implemented", n.Feature)
}
