package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Error struct {
	FieldName    string `json:"fieldName"`
	ErrorMessage string `json:"errorMessage"`
}

type Validator struct {
	v *validator.Validate
}

var StructValidator Validator
var translator ut.Translator

func validateURL(fl validator.FieldLevel) bool {
	URL := fl.Field().String()

	regex, _ := regexp.Compile(`^(https?://)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`)
	result := regex.MatchString(URL)
	return result
}

func Init() {
	StructValidator.v = validator.New()
	// Create a new universal translator instance
	en := en.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")

	// Register the universal translator instance for the validator
	en_translations.RegisterDefaultTranslations(StructValidator.v, translator)

	StructValidator.v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	// Customize the error messages for each validation tag
	StructValidator.v.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	StructValidator.v.RegisterTranslation("email", translator, func(ut ut.Translator) error {
		return ut.Add("email", "{0} is not a valid email address.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
	StructValidator.v.RegisterTranslation("min", translator, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must be at least {1} characters.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})
	StructValidator.v.RegisterTranslation("URL", translator, func(ut ut.Translator) error {
		return ut.Add("URL", "invalid url", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("URL", fe.Field(), fe.Param())
		return t
	})
	StructValidator.v.RegisterValidation("URL", validateURL)
}

// validate
func (v *Validator) Validate(val interface{}) map[string]interface{} {
	errs := make(map[string]interface{})

	if err := StructValidator.v.Struct(val); err != nil {
		// Check if the error is a ValidationErrors slice
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			i := 0
			for _, validationError := range validationErrors {
				errs[validationError.Field()] = Error{
					FieldName:    validationError.Field(),
					ErrorMessage: validationError.Translate(translator),
				}
				i++
			}
			errs["first"] = validationErrors[0].Translate(translator)
			errs["errorsLength"] = i
		}
	}
	return errs
}
