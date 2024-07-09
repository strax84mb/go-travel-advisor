package validation

type Validator[T int] func(T) error

func Validate[T int](
	name string,
	value T,
	ve *ValidationError,
	validators ...Validator[T],
) *ValidationError {
	var (
		pve *ParamValidationError
		err error
	)
	for _, validator := range validators {
		err = validator(value)
		if err != nil {
			if pve == nil {
				pve = NewParamValidationError(name)
			}
			pve.AddError(err)
		}
	}
	if pve != nil {
		if ve == nil {
			ve = &ValidationError{}
		}
		ve.AddError(pve)
	}
	return ve
}
