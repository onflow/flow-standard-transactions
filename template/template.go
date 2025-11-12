package template

import (
	"fmt"
)

type SimpleTemplate struct {
	name        string
	label       Label
	cardinality uint

	initialParameters Parameters
	transactionEdit   func(parameters Parameters) (TransactionEdit, error)
}

var _ Template = (*SimpleTemplate)(nil)

func NewSimpleTemplate(
	name string,
	label Label,
	cardinality uint,
) *SimpleTemplate {
	return &SimpleTemplate{
		name:        name,
		label:       label,
		cardinality: cardinality,
	}
}

func (s *SimpleTemplate) Name() string {
	return s.name
}

func (s *SimpleTemplate) Label() Label {
	return s.label
}

func (s *SimpleTemplate) Cardinality() uint {
	return s.cardinality
}

func (s *SimpleTemplate) WithTransactionEdit(
	transactionEdit func(parameters Parameters) (TransactionEdit, error),
) *SimpleTemplate {
	s.transactionEdit = transactionEdit
	return s
}

func (s *SimpleTemplate) TransactionEditFunc(parameters Parameters) (TransactionEdit, error) {
	if s.transactionEdit == nil {
		return TransactionEdit{}, nil
	}

	return s.transactionEdit(parameters)
}

func (s *SimpleTemplate) InitialParameters() Parameters {
	if s.initialParameters == nil {
		return make(Parameters, s.cardinality)
	}

	// Copy the initial parameters to avoid modifying the original
	cloned := make(Parameters, len(s.initialParameters))
	copy(cloned, s.initialParameters)

	return cloned
}

func (s *SimpleTemplate) WithInitialParameters(
	initialParameters Parameters,
) *SimpleTemplate {
	s.initialParameters = initialParameters
	return s
}

func LoopTemplate(
	n uint64,
	body string,
) string {
	return fmt.Sprintf(`
				var i = 0
				while i < %d {
					i = i + 1
					%s
				}`, n, body)
}
