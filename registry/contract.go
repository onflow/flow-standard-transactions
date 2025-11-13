package registry

import (
	_ "embed"
	"fmt"

	"github.com/onflow/flow-standard-transactions/template"
)

//go:embed contract.cdc
var contract []byte

func simpleContractTemplate(
	name string,
	label template.Label,
	cardinality uint,
	bodyFunc func(parameters template.Parameters, te *template.TransactionEdit) error,
) *template.SimpleTemplate {
	return template.NewSimpleTemplate(
		name,
		label,
		cardinality,
	).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				te := template.TransactionEdit{}
				err := bodyFunc(parameters, &te)
				if err != nil {
					return template.TransactionEdit{}, err
				}
				return te, nil
			},
		)
}

func simpleContractTemplateWithLoop(
	name string,
	label template.Label,
	initialLoopLength uint64,
	body string,
) *template.SimpleTemplate {
	return simpleContractTemplate(
		name,
		label,
		1,
		func(parameters template.Parameters, te *template.TransactionEdit) error {
			te.PrepareBlock = template.LoopTemplate(parameters[0], body)
			return nil
		},
	).WithInitialParameters(template.Parameters{initialLoopLength})
}

var ContractTemplates = []template.Template{
	simpleContractTemplateWithLoop(
		"call empty contract function",
		"CEC",
		1144,
		`TestContract.empty()`,
	),
	simpleContractTemplateWithLoop(
		"emit event",
		"CEE",
		799,
		`TestContract.emitEvent()`,
	),
	simpleContractTemplateWithLoop(
		"mint NFT",
		"CMNFT",
		157,
		`TestContract.mintNFT()`,
	),
	simpleContractTemplate(
		"emit event with string",
		"CEES",
		1,
		func(parameters template.Parameters, te *template.TransactionEdit) error {
			body := fmt.Sprintf(`
					let dict: {String: String} = %s
					TestContract.emitDictEvent(dict)
				`, stringDictOfLen(parameters[0], 50))

			te.PrepareBlock = body
			return nil
		},
	).WithInitialParameters(template.Parameters{960}),
}
