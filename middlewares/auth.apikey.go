package middleware

import (
	"net/http"
	"urlshortner/config/dotenv"
	Logger "urlshortner/utils/logger"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		ApiKey := ctx.GetHeader("x-api-key")
		defer Logger.Log.Sync()
		if ApiKey == "" || ApiKey != dotenv.Global.ApiKey {
			Logger.Log.Error("invalid api key")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid api key",
			})
		}
		ctx.Next()
	})
}
