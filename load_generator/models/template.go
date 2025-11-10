package models

type Template interface {
	Name() string
	Label() Label

	// Cardinality returns the number of parameters that this template has
	Cardinality() uint
}
