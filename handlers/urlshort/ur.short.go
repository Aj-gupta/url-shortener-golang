package shorturlhandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"urlshortner"
	shortUrlCtrl "urlshortner/controllers/urlshort"
	"urlshortner/utils"
	"urlshortner/utils/logger"
	validator "urlshortner/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	service shortUrlCtrl.Service
}

func NewShortUrlHandler(svc shortUrlCtrl.Service) *handler {
	return &handler{service: svc}
}

func (h *handler) GetShortUrl(ctx *gin.Context) (response, headers interface{}, err error) {
	defer logger.Log.Sync()
	url := ctx.Query("url")
	body := &shortUrlCtrl.CreateShortUrl{}
	body.URL = url
	validationErrs := validator.StructValidator.Validate(body)
	if len(validationErrs) > 0 {
		logger.Log.Warn("error while validation request body", zap.Any("data", validationErrs))
		return nil, nil, utils.CustomAPIErr{Code: http.StatusBadRequest, Err: errors.New("invalid request"), Errors: validationErrs}
	}
	res, err := h.service.GetShortUrl(*body)

	if err != nil {
		return nil, nil, err
	}

	genericResponse := &urlshortner.GenericResponse{}
	genericResponse.Data = res
	genericResponse.Message = "successfully get url"

	return genericResponse, nil, nil
}

func (h *handler) CreateShortUrl(ctx *gin.Context) (response, headers interface{}, err error) {

	body := &shortUrlCtrl.CreateShortUrl{}
	requestBody, err := ioutil.ReadAll(ctx.Request.Body)
	defer logger.Log.Sync()

	if err != nil {
		logger.Log.Error("error while decoding request body err=" + err.Error())
		return nil, nil, utils.CustomAPIErr{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}
	err = json.Unmarshal(requestBody, body)
	if err != nil {
		logger.Log.Error("error while unmarshal request body err=" + err.Error())
		return nil, nil, utils.CustomAPIErr{Code: http.StatusUnauthorized, Err: errors.New("invalid request body")}
	}

	validationErrs := validator.StructValidator.Validate(body)
	if len(validationErrs) > 0 {
		logger.Log.Warn("error while validation request body", zap.Any("data", validationErrs))
		return nil, nil, utils.CustomAPIErr{Code: http.StatusBadRequest, Err: errors.New("invalid request body"), Errors: validationErrs}
	}

	res, err := h.service.CreateShortUrl(*body)
	if err != nil {
		return nil, nil, err

	}

	genericResponse := &urlshortner.GenericResponse{}
	genericResponse.Data = res
	genericResponse.Message = "successfully created short url"

	return genericResponse, nil, nil
}
