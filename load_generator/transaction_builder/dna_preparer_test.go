package transaction_builder_test

import (
	"context"
	"testing"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/registry"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/templates"

	"github.com/onflow/flow-go/model/flow"
	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"
)

func Test_TransactionNew(t *testing.T) {
	t.Parallel()

	reg := registry.NewRegistry()
	err := reg.Register(templates.NewSimpleTemplate(
		"T1",
		"T1",
		1,
	).WithTransactionEdit(
		func(parameters models.Parameters) models.TransactionEditFunc {
			return func(
				context models.Context,
				account models.Account,
			) (models.TransactionEdit, error) {
				return models.TransactionEdit{
					PrepareBlock: templates.LoopTemplate(parameters[0], "signer.address"),
				}, nil
			}
		}),
	)
	require.NoError(t, err)

	dna := models.DNA{
		models.DNAElement{
			Label:      "T1",
			Parameters: models.Parameters{5},
		},
	}

	preparer := templates.NewDNAPreparer(zerolog.Nop(), reg)

	c := models.Context{
		ChainID: flow.Emulator,
	}

	body, err := preparer.Prepare(context.Background(), dna, c, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, body)
}
