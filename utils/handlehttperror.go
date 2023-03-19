package utils

import "github.com/gin-gonic/gin"

func HandleHTTPError(err error) (int, gin.H) {
	var statusCode int
	var errors map[string]interface{}
	errorMessage := err.Error()
	errorsLength := 1
	pgErr, isPGErr := HandlePostgresError(err)
	if isPGErr {
		err = pgErr
		errorMessage = err.Error()
	}
	switch e := err.(type) {
	case CustomAPIErr:
		statusCode = e.Status()
		errors = e.ErrorArr()
		if errors["errorsLength"] != nil {
			errorsLength = errors["errorsLength"].(int)
			if errors["errorsLength"] == 1 {
				errorMessage = errors["first"].(string)
			}
		} else {
			errors = make(map[string]interface{})
			errors["errorsLength"] = errorsLength
			errors["first"] = errorMessage
		}

	default:
		errors = make(map[string]interface{})
		errors["errorsLength"] = errorsLength
		errors["first"] = errorMessage
	}
	return statusCode, gin.H{
		"error":   true,
		"message": errorMessage,
		"errors":  errors,
	}
}
