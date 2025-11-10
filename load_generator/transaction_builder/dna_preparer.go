package transaction_builder

import (
	"fmt"
	"strings"

	"github.com/onflow/flow-standard-transactions/load_generator/models"

	"github.com/rs/zerolog"
)

// DNAPreparer
// is responsible for keeping track of which setup functions have been called
type TransactionBuilder struct {
	registry models.Registry
}

func NewTransactionBuilder(
	log zerolog.Logger,
	registry models.Registry,
) *TransactionBuilder {
	return &TransactionBuilder{
		registry: registry,
	}
}

type builtTransaction struct {
	log zerolog.Logger

	transactionBuilder *TransactionBuilder
}

var emptyBody = models.TransactionBody("")

// Prepare prepares everything for sending a transaction with the given a DNA
// - if a global setup is required, it is run once.
// - if an account one time setup is required, it is run once per account.
// - if an account setup is required, it is run every time prepare is called.
//
// It returns the transaction body or an error if something went wrong
func (t *TransactionBuilder) BuildTransaction() (models.TransactionBody, error) {
	preparedDNA := &builtTransaction{
		transactionBuilder: t,
	}

	body, err := preparedDNA.BuildTransactionBody()
	if err != nil {
		// return emptyBody, errors.NewPrepareTransactionFromDNA(err, dna)
	}

	return models.TransactionBody(body), nil
}

func (b *builtTransaction) BuildTransactionBody() (string, error) {
	prepare := strings.Builder{}
	execute := strings.Builder{}
	declare := strings.Builder{}

	// imports need to be added later by the user of this library
	return fmt.Sprintf(
		`
transaction(){
    %s

    prepare(signer: auth(Storage, Contracts, Keys, Inbox, Capabilities) &Account) {
        var f: fun(): Void = fun(){}
%s    }

    execute {
        var f: fun(): Void = fun(){}
%s    }
}`,
		declare.String(),
		prepare.String(),
		execute.String(),
	), nil
}
