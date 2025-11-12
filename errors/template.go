package errors

import (
	"fmt"

	"github.com/onflow/flow-standard-transactions/template"
)

type TemplateAlreadyRegistered struct {
	Label template.Label
}

var _ Error = (*TemplateAlreadyRegistered)(nil)

func NewTemplateAlreadyRegistered(label template.Label) TemplateAlreadyRegistered {
	return TemplateAlreadyRegistered{Label: label}
}

func (t TemplateAlreadyRegistered) Error() string {
	return fmt.Sprintf("template %s already registered", t.Label)
}

type TemplateNotFound struct {
	Label template.Label
}

var _ Error = (*TemplateNotFound)(nil)

func NewTemplateNotFound(label template.Label) TemplateNotFound {
	return TemplateNotFound{Label: label}
}

func (t TemplateNotFound) Error() string {
	return fmt.Sprintf("template %s not found", t.Label)
}

type ModelCalibration struct {
	Label template.Label
	Err   error
}

var _ UnwrappableError = (*ModelCalibration)(nil)

func NewModelCalibration(label template.Label, err error) ModelCalibration {
	return ModelCalibration{Label: label, Err: err}
}

func (m ModelCalibration) Error() string {
	return fmt.Sprintf("model %s calibration failed: %v", m.Label, m.Err)
}

func (m ModelCalibration) Unwrap() error {
	return m.Err
}

type NoUniqueIntensities struct {
	Label template.Label
}

var _ Error = (*NoUniqueIntensities)(nil)

func NewNoUniqueIntensities(label template.Label) *NoUniqueIntensities {
	return &NoUniqueIntensities{Label: label}
}

func (n *NoUniqueIntensities) Error() string {
	return fmt.Sprintf("model %s has no unique intensities", n.Label)
}
