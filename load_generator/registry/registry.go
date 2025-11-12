package registry

import (
	"github.com/onflow/flow-standard-transactions/load_generator/errors"
	"github.com/onflow/flow-standard-transactions/load_generator/models"
)

type Registry struct {
	templates map[models.Label]models.Template
}

var _ models.Registry = &Registry{}

func NewRegistry() Registry {
	return Registry{
		templates: make(map[models.Label]models.Template),
	}
}

func (r Registry) Clone() Registry {
	templ := make(map[models.Label]models.Template)
	for label, template := range r.templates {
		templ[label] = template
	}
	return Registry{
		templates: templ,
	}
}

func (r Registry) Register(template models.Template) error {
	if _, ok := r.templates[template.Label()]; ok {
		return errors.NewTemplateAlreadyRegistered(template.Label())
	}

	r.templates[template.Label()] = template
	return nil
}

func (r Registry) RegisterAll(template []models.Template) error {
	for _, templ := range template {
		if err := r.Register(templ); err != nil {
			return err
		}
	}
	return nil
}

func (r Registry) Unregister(labels ...models.Label) Registry {
	for _, label := range labels {
		delete(r.templates, label)
	}
	return r
}

func (r Registry) Get(label models.Label) (models.Template, error) {
	templ, ok := r.templates[label]
	if !ok {
		return nil, errors.NewTemplateNotFound(label)
	}

	return templ, nil
}

func (r Registry) AllLabels() []models.Label {
	labels := make([]models.Label, 0, len(r.templates))
	for label := range r.templates {
		labels = append(labels, label)
	}
	return labels
}

var Global Registry

func init() {
	Global = NewRegistry()
	err := Global.RegisterAll(simpleTemplates)
	if err != nil {
		panic(err)
	}
	err = Global.RegisterAll(contractTemplates)
	if err != nil {
		panic(err)
	}

	err = Global.RegisterAll(scheduledTransactions)
	if err != nil {
		panic(err)
	}
}

type InitialParametersData struct {
	Templates map[models.Label]TemplateInitialParameters `json:"templates"`
}

type TemplateInitialParameters struct {
	InitialParameters []uint64 `json:"parameters"`
}
