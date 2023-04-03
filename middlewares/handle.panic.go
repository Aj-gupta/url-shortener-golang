package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"urlshortner/utils"
	"urlshortner/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				defer logger.Log.Sync()
				logger.Log.Error("panic: ", zap.Any("error: ", r))
				errMsg := fmt.Sprintf("panic: %v", r)

				err := errors.New(errMsg)
				_, body := utils.HandleHTTPError(err)

				c.AbortWithStatusJSON(http.StatusInternalServerError, body)
			}
		}()
		c.Next()
	}
}
