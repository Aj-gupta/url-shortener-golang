package utils

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomHandler func(context *gin.Context) (interface{}, interface{}, error)

func HandlerEncoder(handler CustomHandler) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		response, headers, err := handler(c)

		if c.Writer.Written() {
			return
		}
		if err != nil {
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

			c.AbortWithStatusJSON(statusCode, gin.H{
				"error":   true,
				"message": errorMessage,
				"errors":  errors,
			})
			return
		}
		if headers != nil {
			for header_key, header_value := range headers.(map[string]string) {
				c.Header(header_key, header_value)
			}
		}
		// fmt.Println("encoder")
		// reponse type
		switch response.(type) {
		case *bytes.Buffer:
			c.Writer.Write(response.(*bytes.Buffer).Bytes())
		default:
			c.JSON(http.StatusOK, response)
		}
	}
	return gin.HandlerFunc(fn)
}
