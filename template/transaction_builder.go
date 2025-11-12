package template

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type TransactionBuilder struct {
	registry Registry
}

func NewTransactionBuilder(
	log zerolog.Logger,
	registry Registry,
) *TransactionBuilder {
	return &TransactionBuilder{
		registry: registry,
	}
}

type builtTransaction struct {
	log zerolog.Logger

	transactionBuilder *TransactionBuilder
	transactionEdits   []TransactionEdit
}

func (t *TransactionBuilder) BuildTransaction() TransactionBody {
	builtTransaction := &builtTransaction{
		transactionBuilder: t,
	}

	body := builtTransaction.BuildTransactionBody()

	return TransactionBody(body)
}

func (b *builtTransaction) BuildTransactionBody() string {
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
			prepare.WriteString(TrimAndReplaceIndentation(edit.PrepareBlock, 12))
			prepare.WriteString("        }\n        f()\n")
		}

		if edit.ExecuteBlock != "" {
			execute.WriteString("        f = fun() {\n")
			execute.WriteString(TrimAndReplaceIndentation(edit.ExecuteBlock, 12))
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
	)
}
