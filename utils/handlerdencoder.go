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
			statusCode, errBody := HandleHTTPError(err)
			c.AbortWithStatusJSON(statusCode, errBody)
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
