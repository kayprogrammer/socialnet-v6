package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	customValidator *validator.Validate
	translator      ut.Translator
)

func (e *ErrorResponse) Error() string {
	return e.Message
}

// Initialize the custom validator and translator
func init() {
	customValidator = validator.New()
	en := en.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")

	// Register Custom Validators
	customValidator.RegisterValidation("date", DateValidator)
	customValidator.RegisterValidation("reaction_type_validator", ReactionTypeValidator)
	customValidator.RegisterValidation("file_type_validator", FileTypeValidator)
	customValidator.RegisterValidation("usernames_to_update_validator", DistinctField)

	customValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

}

// Register translations
func registerTranslations(param string) {
	// Register custom error messages for each validation tag
	registerTranslation := func(tag string, translation string, translator ut.Translator) {
		customValidator.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
			return ut.Add(tag, translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
	}

	registerTranslation("gt", "Value is too small!", translator)
	registerTranslation("required", "This field is required.", translator)
	registerTranslation("required_if", "This field is required.", translator)
	registerTranslation("required_without", "This field is required.", translator)

	minErrMsg := fmt.Sprintf("%s characters min", param)
	registerTranslation("min", minErrMsg, translator)
	maxErrMsg := fmt.Sprintf("%s characters max", param)
	registerTranslation("max", maxErrMsg, translator)
	registerTranslation("email", "Invalid Email", translator)
	eqErrMsg := fmt.Sprintf("Must be %s", param)
	registerTranslation("eq", eqErrMsg, translator)
}

// CustomValidator is a custom validator that uses "github.com/go-playground/validator/v10"
type CustomValidator struct{}

// Validate performs the validation of the given struct
func (cv *CustomValidator) Validate(i interface{}) *ErrorResponse {
	if err := customValidator.Struct(i); err != nil {
		err := err.(validator.ValidationErrors)
		return cv.translateValidationErrors(err)
	}
	return nil
}

// translateValidationErrors translates the validation errors to custom errors
func (cv *CustomValidator) translateValidationErrors(errs validator.ValidationErrors) *ErrorResponse {
	errData := make(map[string]string)
	for _, err := range errs {
		errParam := err.Param()
		registerTranslations(errParam)
		errMsg := err.Translate(translator)
		errField := err.Field()

		errData[errField] = errMsg
	}
	errResp := RequestErr(ERR_INVALID_ENTRY, "Invalid Entry", errData)
	return &errResp
}

// New creates a new instance of CustomValidator
func Validator() *CustomValidator {
	return &CustomValidator{}
}
