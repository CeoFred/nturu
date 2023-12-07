package validators

import (
	"github.com/nturu/microservice-template/internal/handlers"
	"github.com/nturu/microservice-template/internal/helpers"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
)

var Validator *validator.Validate
var Trans ut.Translator

func init() {
	Validator = validator.New()
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	Trans, _ = uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(Validator, Trans)

	// Custom error messages
	_ = Validator.RegisterTranslation("required", Trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = Validator.RegisterTranslation("len", Trans, func(ut ut.Translator) error {
		return ut.Add("len", "{0} must be exactly {1} characters long", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("len", fe.Field(), fe.Param())
		return t
	})

	_ = Validator.RegisterTranslation("email", Trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} is not a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
}
func ValidateAccountResetScheme(c *fiber.Ctx) error {
	body := new(helpers.AccountReset)
	err := c.BodyParser(&body)
	if err != nil {
		return helpers.Dispatch400Error(c, "invalid payload", nil)
	}

	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidateOTPVerifySchema(c *fiber.Ctx) error {
	body := new(helpers.OtpVerify)
	err := c.BodyParser(&body)
	if err != nil {
		return helpers.Dispatch400Error(c, "invalid payload", nil)
	}

	err = Validator.Struct(body)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidateRegisterUserSchema(c *fiber.Ctx) error {
	body := new(helpers.InputCreateUser)
	err := c.BodyParser(&body)
	if err != nil {
		return helpers.Dispatch400Error(c, "invalid payload", nil)
	}

	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidateLoginUser(c *fiber.Ctx) error {
	body := new(helpers.AuthenticateUser)

	err := c.BodyParser(&body)
	if err != nil {
		return helpers.Dispatch400Error(c, "invalid payload", nil)
	}
	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidateCompleteUserProfile(c *fiber.Ctx) error {
	body := new(helpers.UpdateAccountInformation)

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidateUpdateUserProfile(c *fiber.Ctx) error {
	body := new(handlers.UpdateUserProfileInput)

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}

func ValidatePasswordReset(c *fiber.Ctx) error {
	body := new(handlers.ResetPassword)

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	err = Validator.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var validationMessages []string
		for _, e := range validationErrors {
			validationMessages = append(validationMessages, e.Translate(Trans))
		}
		return helpers.Dispatch400Error(c, "validation failed", validationMessages)
	}
	return c.Next()
}
