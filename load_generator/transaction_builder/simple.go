package transaction_builder

import (
	"fmt"

	"github.com/onflow/flow-standard-transactions/load_generator/models"
)

type SimpleTemplate struct {
	name        string
	label       models.Label
	cardinality uint
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
