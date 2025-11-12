package transaction_builder

import (
	"fmt"
	"strings"

	"github.com/onflow/flow-standard-transactions/load_generator/models"

	"github.com/rs/zerolog"
)

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
	transactionEdits   []models.TransactionEdit
}

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

	for _, edit := range b.transactionEdits {
		if edit.FieldDeclarations != "" {
			declare.WriteString(strings.Trim(edit.FieldDeclarations, "\n\r"))
			declare.WriteString("\n")
		}

		if edit.PrepareBlock != "" {
			prepare.WriteString("        f = fun() {\n")
			prepare.WriteString(models.TrimAndReplaceIndentation(edit.PrepareBlock, 12))
			prepare.WriteString("        }\n        f()\n")
		}

		if edit.ExecuteBlock != "" {
			execute.WriteString("        f = fun() {\n")
			execute.WriteString(models.TrimAndReplaceIndentation(edit.ExecuteBlock, 12))
			execute.WriteString("        }\n        f()\n")
		}

	}

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
