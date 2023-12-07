package helpers

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateBody(body interface{}) []string {
	Validator := validator.New()
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	Trans, _ := uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(Validator, Trans)

	err := Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return validationMessages
	}
	return nil
}
