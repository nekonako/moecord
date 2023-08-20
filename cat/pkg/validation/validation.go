package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type responseError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func IsValidationError(e error) (bool, []responseError) {
	ve, ok := e.(validator.ValidationErrors)
	if !ok {
		return false, nil
	}
	result := make([]responseError, len(ve))
	for i, v := range ve {
		m := []string{}
		if v.Tag() != "" {
			m = append(m, v.Tag())
		}
		if v.Param() != "" {
			m = append(m, v.Param())
		}
		result[i] = responseError{
			Field:   v.Field(),
			Message: strings.Join(m, " "),
		}
	}
	return true, result
}

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
