package registry

import (
	_ "embed"
	"fmt"

	"github.com/onflow/flow-standard-transactions/load_generator/models"
	"github.com/onflow/flow-standard-transactions/load_generator/transaction_builder"
)

//go:embed contract.cdc
var contract []byte

func simpleContractTemplate(
	name string,
	label models.Label,
	cardinality uint,
	bodyFunc func(parameters models.Parameters, te *models.TransactionEdit) error,
) *transaction_builder.SimpleTemplate {
	return transaction_builder.NewSimpleTemplate(
		name,
		label,
		cardinality,
	).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEdit {
				te := models.TransactionEdit{}
				err := bodyFunc(parameters, &te)
				if err != nil {
					return models.TransactionEdit{}
				}
				return te
			},
		)
}

func simpleContractTemplateWithLoop(
	name string,
	label models.Label,
	initialLoopLength uint64,
	body string,
) *transaction_builder.SimpleTemplate {
	return simpleContractTemplate(
		name,
		label,
		1,
		func(parameters models.Parameters, te *models.TransactionEdit) error {
			te.PrepareBlock = transaction_builder.LoopTemplate(parameters[0], body)
			return nil
		},
	).WithInitialParameters(models.Parameters{initialLoopLength})
}

var contractTemplates = []models.Template{
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
		func(parameters models.Parameters, te *models.TransactionEdit) error {
			body := fmt.Sprintf(`
					let dict: {String: String} = %s
					TestContract.emitDictEvent(dict)
				`, stringDictOfLen(parameters[0], 50))

			te.PrepareBlock = body
			return nil
		},
	).WithInitialParameters(models.Parameters{960}),
}
