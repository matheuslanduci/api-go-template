package services

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	Validator *validator.Validate
}

func NewValidator() *XValidator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &XValidator{
		Validator: validate,
	}
}

type ErrorResponse struct {
	Property string `json:"property"`
	Rule     string `json:"rule"`
}

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := v.Validator.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.Property = err.Field()

			if err.Param() != "" {
				elem.Rule = fmt.Sprintf("%s:%s", err.Tag(), err.Param())
			} else {
				elem.Rule = err.Tag()
			}

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
