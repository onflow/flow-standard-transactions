package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/flow"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/templates"
)

const scheduleTemplate = `
	if !signer.storage.check<@TestContract.Handler>(from: TestContract.HandlerStoragePath) {
		let handler <- TestContract.createHandler()

		signer.storage.save(<-handler, to: TestContract.HandlerStoragePath)
		signer.capabilities.storage.issue<auth(FlowTransactionScheduler.Execute) &{FlowTransactionScheduler.TransactionHandler}>(TestContract.HandlerStoragePath)
	}

	let handlerCap = signer.capabilities.storage
						.getControllers(forPath: TestContract.HandlerStoragePath)[0]
						.capability as! Capability<auth(FlowTransactionScheduler.Execute) &{FlowTransactionScheduler.TransactionHandler}>

	let vault = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)
		?? panic("Could not borrow FlowToken vault")

	%s

	let scheduledTransaction <- FlowTransactionScheduler.schedule(
		handlerCap: handlerCap,
		data: data,
		timestamp: timestamp,
		priority: priority,
		executionEffort: effort,
		fees: <-fees
	)
	destroy scheduledTransaction
`

func simpleScheduledTransactionTemplateWithLoop(
	name string,
	label models.Label,
	initialParams models.Parameters,
	cardinality uint,
	body func(models.Parameters) string,
) *templates.SimpleTemplate {
	var contractAddress *flow.Address

	return templates.NewSimpleTemplate(
		name,
		label,
		cardinality,
	).
		WithInitialParameters(initialParams).
		WithTransactionEdit(func(parameters models.Parameters) models.TransactionEditFunc {
			return func(
				context models.Context,
				account models.Account,
			) (models.TransactionEdit, error) {
				sc := systemcontracts.SystemContractsForChain(context.ChainID)

				return models.TransactionEdit{
					Imports: map[string]flow.Address{
						sc.FlowToken.Name:             sc.FlowToken.Address,
						sc.FungibleToken.Name:         sc.FungibleToken.Address,
						sc.FlowCallbackScheduler.Name: sc.FlowCallbackScheduler.Address,
						"TestContract":                *contractAddress,
					},
					PrepareBlock: templates.LoopTemplate(parameters[0], body(parameters)),
				}, nil
			}
		}).
		WithGlobalSetup(func(
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
		})
}

var scheduledTransactions = []models.Template{
	simpleScheduledTransactionTemplateWithLoop(
		"scheduled transaction and execute",
		"ST",
		models.Parameters{1},
		1,
		func(params models.Parameters) string {
			return fmt.Sprintf(scheduleTemplate, `
				let fees <- vault.withdraw(amount: 0.003) as! @FlowToken.Vault
				let timestamp = getCurrentBlock().timestamp + 120.0 // 2 minutes in future
				let effort: UInt64 = 100
				let priority = FlowTransactionScheduler.Priority.High
				let data: UInt64 = 0
			`)
		},
	),
	// first param is loop length -> how many scheduled transactions there are
	// second param is data size -> how much data is in each scheduled transaction (100 means 10kb string)
	simpleScheduledTransactionTemplateWithLoop(
		"scheduled transaction and execute with large data (100KB)",
		"STLD",
		models.Parameters{1, 1},
		2,
		func(params models.Parameters) string {
			return fmt.Sprintf(scheduleTemplate, fmt.Sprintf(`
				let fees <- vault.withdraw(amount: 0.11) as! @FlowToken.Vault
				let timestamp = getCurrentBlock().timestamp + 120.0 // 2 minutes in future
				let effort: UInt64 = 100
				let priority = FlowTransactionScheduler.Priority.High
				let data = "%s"
			`, strings.Repeat("A", int(100*params[1])))) // inject N KB of data
		},
	),
	// first param is loop length -> how many scheduled transactions there are
	// second param is data size -> how much data is in each scheduled transaction ( 100 means array of 100 elements)
	simpleScheduledTransactionTemplateWithLoop(
		"scheduled transaction and execute with large array (10k items)",
		"STLA",
		models.Parameters{1, 1},
		2,
		func(params models.Parameters) string {
			return fmt.Sprintf(scheduleTemplate, fmt.Sprintf(`
				let largeArray: [Int] = []
				while largeArray.length < %d {
					largeArray.append(1)
				}
	
				let fees <- vault.withdraw(amount: 0.01) as! @FlowToken.Vault
				let timestamp = getCurrentBlock().timestamp + 120.0 // 2 minutes in future
				let effort: UInt64 = 100
				let priority = FlowTransactionScheduler.Priority.High
				let data = largeArray
			`, params[1]))
		},
	),
}
