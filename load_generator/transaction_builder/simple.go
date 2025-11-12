package transaction_builder

import (
	"fmt"

	"github.com/onflow/flow-standard-transactions/load_generator/models"
)

type SimpleTemplate struct {
	name        string
	label       models.Label
	cardinality uint

	initialParameters models.Parameters
	transactionEdit   func(parameters models.Parameters) models.TransactionEdit
}

var _ models.Template = (*SimpleTemplate)(nil)

func NewSimpleTemplate(
	name string,
	label models.Label,
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

func (s *SimpleTemplate) Label() models.Label {
	return s.label
}

func (s *SimpleTemplate) Cardinality() uint {
	return s.cardinality
}

func (s *SimpleTemplate) WithTransactionEdit(
	transactionEdit func(parameters models.Parameters) models.TransactionEdit,
) *SimpleTemplate {
	s.transactionEdit = transactionEdit
	return s
}

func (s *SimpleTemplate) TransactionEdit(parameters models.Parameters) models.TransactionEdit {
	if s.transactionEdit == nil {
		return models.TransactionEdit{}
	}

	return s.transactionEdit(parameters)
}

func (s *SimpleTemplate) InitialParameters() models.Parameters {
	if s.initialParameters == nil {
		return make(models.Parameters, s.cardinality)
	}

	// Copy the initial parameters to avoid modifying the original
	cloned := make(models.Parameters, len(s.initialParameters))
	copy(cloned, s.initialParameters)

	return cloned
}

func (s *SimpleTemplate) WithInitialParameters(
	initialParameters models.Parameters,
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
