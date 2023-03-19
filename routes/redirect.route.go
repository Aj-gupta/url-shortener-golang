package routes

import (
	shortUrlCtrl "urlshortner/controllers/urlshort"
	redirectHandler "urlshortner/handlers/redirect"
	"urlshortner/utils"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func InitRedirectRoute(db *bun.DB, route *gin.Engine) {

	ShortUrlRepo := shortUrlCtrl.NewShortUrlRepo(db)
	ShortUrlSvc := shortUrlCtrl.NewShortUrlSvc(ShortUrlRepo)
	RedirectHandler := redirectHandler.NewRedirectHandler(ShortUrlSvc)

	redirectGroup := route.Group("/:shortUrl")
	redirectGroup.GET("", utils.HandlerEncoder(RedirectHandler.Redirect))
}
