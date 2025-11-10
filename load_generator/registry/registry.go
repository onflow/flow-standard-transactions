package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/errors"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/templates"
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

func (r Registry) ValidateDNA(dna models.DNA) error {
	for _, element := range dna {
		template, err := r.Get(element.Label)
		if err != nil {
			return fmt.Errorf("template not found: %w", err)
		}
		if template.Cardinality() != uint(len(element.Parameters)) {
			return fmt.Errorf("template %s has cardinality %d, but DNA has %d parameters", element.Label, template.Cardinality(), len(element.Parameters))
		}
	}
	return nil
}

func (r Registry) AllLabels() []models.Label {
	labels := make([]models.Label, 0, len(r.templates))
	for label := range r.templates {
		labels = append(labels, label)
	}
	return labels
}

func (r Registry) LoadInitialParametersFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var initialParameters InitialParametersData
	err = json.Unmarshal(file, &initialParameters)
	if err != nil {
		return err
	}

	for label, parameters := range initialParameters.Templates {
		t, ok := r.templates[label]
		if !ok {
			continue
		}

		switch tt := t.(type) {
		case *templates.SavedInitialParameters:
			t = templates.NewSavedInitialParameters(tt.Template, parameters.InitialParameters)
		default:
			t = templates.NewSavedInitialParameters(tt, parameters.InitialParameters)
		}
		r.templates[label] = t
	}
	return nil
}

func SaveInitialParametersToFile(path string, data InitialParametersData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0o644)
	if err != nil {
		return err
	}
	return nil
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
