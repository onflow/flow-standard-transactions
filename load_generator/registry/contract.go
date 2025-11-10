package registry

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/templates"

	"github.com/onflow/flow-go/fvm/blueprints"
	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/flow"
)

//go:embed contract.cdc
var contract []byte

func setupTestContract(
	ctx context.Context,
	c models.Context,
	setup models.ChainInteraction,
) (flow.Address, error) {
	// replace resolved imports in contract code
	sc := systemcontracts.SystemContractsForChain(c.ChainID)
	script := strings.ReplaceAll(
		string(contract),
		`import "FlowTransactionScheduler"`,
		fmt.Sprintf("import FlowTransactionScheduler from %s\n", sc.FlowCallbackScheduler.Address.HexWithPrefix()),
	)

	return setupContract(ctx, c, setup, []byte(script), "TestContract")
}

func setupContract(
	ctx context.Context,
	c models.Context,
	setup models.ChainInteraction,
	script []byte,
	name string,
) (flow.Address, error) {
	referenceBlock, err := setup.ReferenceBlock()
	if err != nil {
		return flow.EmptyAddress, err
	}

	// don't return this account
	account, err := setup.Borrow()
	if err != nil {
		return flow.EmptyAddress, err
	}

	deployTx := blueprints.DeployContractTransaction(
		account.Address(),
		script,
		name,
	).
		SetComputeLimit(9999).
		SetReferenceBlockID(referenceBlock).
		SetPayer(account.Address())

	err = account.SetAsProposer(deployTx)
	if err != nil {
		return flow.EmptyAddress, err
	}
	err = account.SignEnvelope(deployTx)
	if err != nil {
		return flow.EmptyAddress, err
	}

	txBody, err := deployTx.Build()
	if err != nil {
		return flow.EmptyAddress, err
	}

	r, err := setup.Send(ctx, txBody)
	if err != nil {
		return flow.EmptyAddress, err
	}
	result := <-r
	account.IncrementSequenceNumber()
	if result.Error != nil {
		return flow.EmptyAddress, result.Error
	}
	return account.Address(), nil
}

func simpleContractTemplate(
	name string,
	label models.Label,
	cardinality uint,
	bodyFunc func(parameters models.Parameters, contractAddress flow.Address, te *models.TransactionEdit) error,
) *templates.SimpleTemplate {
	var contractAddress *flow.Address

	return templates.NewSimpleTemplate(
		name,
		label,
		cardinality,
	).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(
					context models.Context,
					account models.Account,
				) (models.TransactionEdit, error) {
					if contractAddress == nil {
						return models.TransactionEdit{}, fmt.Errorf("contract address not set yet")
					}

					sc := systemcontracts.SystemContractsForChain(context.ChainID)

					te := models.TransactionEdit{
						Imports: map[string]flow.Address{
							sc.FlowToken.Name:     sc.FlowToken.Address,
							sc.FungibleToken.Name: sc.FungibleToken.Address,
							"TestContract":        *contractAddress,
						},
					}

					err := bodyFunc(parameters, *contractAddress, &te)
					if err != nil {
						return models.TransactionEdit{}, err
					}

					return te, nil
				}
			},
		).WithGlobalSetup(
		func(
			ctx context.Context,
			c models.Context,
			setup models.ChainInteraction,
		) error {
			address, err := setupTestContract(ctx, c, setup)
			if err != nil {
				return err
			}
			contractAddress = &address

			return nil
		},
	)
}

func simpleContractTemplateWithLoop(
	name string,
	label models.Label,
	initialLoopLength uint64,
	body string,
) *templates.SimpleTemplate {
	return simpleContractTemplate(
		name,
		label,
		1,
		func(parameters models.Parameters, contractAddress flow.Address, te *models.TransactionEdit) error {
			te.PrepareBlock = templates.LoopTemplate(parameters[0], body)
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
		func(parameters models.Parameters, contractAddress flow.Address, te *models.TransactionEdit) error {
			body := fmt.Sprintf(`
					let dict: {String: String} = %s
					TestContract.emitDictEvent(dict)
				`, stringDictOfLen(parameters[0], 50))

			te.PrepareBlock = body
			return nil
		},
	).WithInitialParameters(models.Parameters{960}),
}
