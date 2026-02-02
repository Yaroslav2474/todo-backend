package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	validate := validator.New()

	_ = validate.RegisterValidation("notfuture", validateNotFutureDate)
	_ = validate.RegisterValidation("notpast", validateNotPastDate)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: validate}

}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (cv *CustomValidator) ValidateField(field interface{}, tag string) error {
	return cv.validator.Var(field, tag)
}

func (cv *CustomValidator) GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			errors = append(errors, formatFieldError(ve))
		}
	}

	return errors
}

func formatFieldError(ve validator.FieldError) string {
	field := ve.Field()
	if ve.Param() != "" {
		return fmt.Sprintf("%s: %s (значение: %v, ожидается: %s)",
			field, ve.Tag(), ve.Value(), ve.Param())
	}
	return fmt.Sprintf("%s: %s (значение: %v)", field, ve.Tag(), ve.Value())
}

func validateNotFutureDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		return !date.After(time.Now())
	}
	return false
}

func validateNotPastDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		return !date.Before(time.Now())
	}
	return false
}
