package validation

import (
	"errors"
	"fmt"
	"strings"
)

type Destroyable interface {
	Destroy()
}

type ValidationError struct {
	errs []*ParamValidationError
}

func (ve *ValidationError) AddError(err *ParamValidationError) {
	ve.errs = append(ve.errs, err)
}

func (ve *ValidationError) Error() string {
	texts := make([]string, len(ve.errs))
	for i, err := range ve.errs {
		texts[i] = err.Error()
	}
	return "Validation errors:" + strings.Join(texts, "\n")
}

func (ve *ValidationError) Destroy() {
	for i, err := range ve.errs {
		err.Destroy()
		ve.errs[i] = nil
	}
}

type ParamValidationError struct {
	name   string
	errors []error
}

func NewParamValidationError(name string) *ParamValidationError {
	return &ParamValidationError{
		name: name,
	}
}

func (pve *ParamValidationError) AddError(err error) {
	pve.errors = append(pve.errors, err)
}

func (pve *ParamValidationError) AddErrorByMessage(msg string) {
	pve.errors = append(pve.errors, errors.New(msg))
}

func (pve *ParamValidationError) Error() string {
	texts := make([]string, len(pve.errors))
	for i, err := range pve.errors {
		texts[i] = err.Error()
	}
	return fmt.Sprintf("%s : \n\t%s", pve.name, strings.Join(texts, "\n\t"))
}

func (pve *ParamValidationError) Destroy() {
	for i := range pve.errors {
		pve.errors[i] = nil
	}
}
