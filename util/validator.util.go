package util

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	Validator *validator.Validate
	once      sync.Once
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func BodyValidator(c *fiber.Ctx, body interface{}) (bool, string, string) {
	initValidator()

	if err := c.BodyParser(body); err != nil {
		return true, "BODY_PARSER_ERROR", "invalid request body"
	}

	// Validation
	validationErrors := []ErrorResponse{}
	errs := Validator.Struct(body)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0)
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return true, "VALIDATION_ERROR", strings.Join(errMsgs, " and ")
	}

	return false, "", ""
}

func initValidator() {
	once.Do(func() {
		Validator = validator.New()
		Validator.RegisterValidation("array-required", validateArrayRequired)
	})
}

func validateArrayRequired(fl validator.FieldLevel) bool {
	return len(fl.Field().Interface().([]any)) > 0
}
