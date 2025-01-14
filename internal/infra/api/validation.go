package api

import (
	"encoding/json"
	"errors"
	"fullcycle-auction_go/configuration/rest_err"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validatorEn "github.com/go-playground/validator/v10/translations/en"
)

var translator ut.Translator

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		englishTranslator := en.New()
		enTransl := ut.New(englishTranslator, englishTranslator)
		translator, _ = enTransl.GetTranslator("en")
		_ = validatorEn.RegisterDefaultTranslations(value, translator)
	}
}

func ValidateErr(validationErr error) *rest_err.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidation validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return rest_err.NewNotFoundError("Invalid type error")
	}
	
	if errors.As(validationErr, &jsonValidation) {
		var errorCauses []rest_err.Causes
		for _, e := range validationErr.(validator.ValidationErrors) {
			errorCauses = append(errorCauses, rest_err.Causes{
				Field:   e.Field(),
				Message: e.Translate(translator),
			})
		}
		return rest_err.NewBadRequestError("Invalid field values", errorCauses...)
	}

	return rest_err.NewBadRequestError("Error trying to convert fields")
}
