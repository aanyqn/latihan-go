package helper

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"unicode"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	_ = v.RegisterValidation("nospace", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), " \t\n\r")
	})
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) &&
				r != '.' && r != '_' {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		return passwordStrength(fl.Field().String()) == ""
	})
	return v
}

func ValidateStruct(s any) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var invalid *validator.InvalidValidationError
	if errors.As(err, &invalid) {
		return map[string]string{"_": "validated object isn't valid"}
	}
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return map[string]string{"_": "validation failed"}
	}
	result := make(map[string]string, len(fieldErrors))
	for _, fe := range fieldErrors {
		if _, exists := result[fe.Field()]; !exists {
			result[fe.Field()] = messageFor(fe)
		}
	}
	return result
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "must be filled"
	case "email":
		return "invalid email format"
	case "min":
		if fe.Kind() == reflect.String {
			return "minimum " + fe.Param() + " character"
		}
		return "minimum value " + fe.Param()
	case "max":
		if fe.Kind() == reflect.String {
			return "maximum " + fe.Param() + " character"
		}
		return "maximal value " + fe.Param()
	case "alphanum":
		return "only letters and digits"
	case "nospace":
		return "can't contains space"
	case "username":
		return "only letters, digits, dots and underscore"
	case "strongpassword":
		if value, ok := fe.Value().(string); ok {
			return passwordStrength(value)
		}
		return "password too weak"
	case "oneof":
		return "must be one of: " +
			strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "not valid " + fe.Tag()
	}
}

const minPasswordLength = 8

func passwordStrength(password string) string {
	if len(password) < minPasswordLength {
		return "8 character minimum"
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "must contain letters and digit"
	}
	weak := map[string]bool{
		"password1": true, "12345678": true, "qwerty123": true,
		"admin123": true, "password123": true,
	}
	if weak[strings.ToLower(password)] {
		return "password to weak"
	}

	return ""
}
