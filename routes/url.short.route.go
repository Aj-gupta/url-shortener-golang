package routes

import (
	shortUrlCtrl "urlshortner/controllers/urlshort"
	shortUrlHandler "urlshortner/handlers/urlshort"
	middleware "urlshortner/middlewares"
	"urlshortner/utils"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func InitShortUrlRoute(db *bun.DB, route *gin.Engine) {

	ShortUrlRepo := shortUrlCtrl.NewShortUrlRepo(db)
	ShortUrlSvc := shortUrlCtrl.NewShortUrlSvc(ShortUrlRepo)
	ShortUrlHandler := shortUrlHandler.NewShortUrlHandler(ShortUrlSvc)

	groupRoute := route.Group("api/v1/url")
	groupRoute.POST("/short", middleware.Auth(), utils.HandlerEncoder(ShortUrlHandler.CreateShortUrl))
	groupRoute.GET("/short", middleware.Auth(), utils.HandlerEncoder(ShortUrlHandler.GetShortUrl))

}
