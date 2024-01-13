package utils

import (
	"fmt"
	"urlshortner/utils/validator"

	"github.com/gin-gonic/gin"
)

func HandleHTTPError(err error) (int, gin.H) {
	var statusCode int
	var errors map[string]interface{}
	errorMessage := validator.Error{FieldName: "", ErrorMessage: err.Error()}
	errorsLength := 1

	switch e := err.(type) {
	case CustomAPIErr:
		statusCode = e.Status()
		errors = e.ErrorArr()

		errorsMeta := map[string]interface{}{}

		if errors["meta"] != nil {
			errorsMeta = errors["meta"].(map[string]interface{})
		}
		if errors["meta"] != nil && errorsMeta["errorsLength"] != nil {
			errorsLength = errorsMeta["errorsLength"].(int)
			if errorsMeta["errorsLength"] == 1 {
				errorMessage = errorsMeta["first"].(validator.Error)
			}
		} else {
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<dfdf>>>>>>>>>>>>>>>>>>>>>>>")
			errors = make(map[string]interface{})
			errors["meta"] = map[string]interface{}{
				"errorsLength": errorsLength,
				"first":        errorMessage,
			}
		}

	default:
		errors = make(map[string]interface{})
		errors["meta"] = map[string]interface{}{
			"errorsLength": errorsLength,
			"first":        errorMessage,
		}
	}
	return statusCode, gin.H{
		"error":   true,
		"message": errorMessage.ErrorMessage,
		"errors":  errors,
	}
}
