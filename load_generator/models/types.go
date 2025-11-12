package models

type Label = string

type Parameters = []uint64

type TransactionBody string

type TransactionEdit struct {
	PrepareBlock      string
	ExecuteBlock      string
	FieldDeclarations string
}
