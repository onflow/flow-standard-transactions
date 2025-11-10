package models

type Registry interface {
	Get(label Label) (Template, error)
	AllLabels() []Label
}
