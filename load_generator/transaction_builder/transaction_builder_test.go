package transaction_builder_test

import (
	"testing"

	"github.com/onflow/flow-standard-transactions/load_generator/registry"
	"github.com/onflow/flow-standard-transactions/load_generator/transaction_builder"
	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"
)

func Test_TransactionNew(t *testing.T) {
	t.Parallel()

	reg := registry.NewRegistry()
	err := reg.Register(transaction_builder.NewSimpleTemplate(
		"T1",
		"T1",
		1,
	))
	require.NoError(t, err)

	body, err := transaction_builder.NewTransactionBuilder(zerolog.Nop(), reg).BuildTransaction()
	require.NoError(t, err)
	require.NotEmpty(t, body)
}
