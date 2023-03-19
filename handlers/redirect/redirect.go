package redirecthanlder

import (
	"net/http"
	shortUrlCtrl "urlshortner/controllers/urlshort"

	"github.com/gin-gonic/gin"
)

type handler struct {
	shortUrlSvc shortUrlCtrl.Service
}

func NewRedirectHandler(shortUrlSvc shortUrlCtrl.Service) *handler {
	return &handler{shortUrlSvc}
}

func (h *handler) Redirect(ctx *gin.Context) (response, headers interface{}, err error) {
	// TO-DO: add caching headers if url not found
	shortUrl := ctx.Param("shortUrl")
	// c.Header("Vary", "Accept-Encoding")
	// c.Header("Cache-Control", "public, max-age=2592000")
	// if we get to this point we should not let the client cache
	res, err := h.shortUrlSvc.GetLongUrl(shortUrl)
	if err != nil {
		return nil, nil, err
	}
	ctx.Header("Cache-Control", "no-cache, no-store")
	ctx.Redirect(http.StatusTemporaryRedirect, res)
	return nil, nil, nil
}
