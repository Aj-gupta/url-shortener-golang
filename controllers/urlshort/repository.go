package shorturl

import (
	"context"
	"errors"
	"net/http"
	"urlshortner/config/dotenv"
	"urlshortner/models"
	"urlshortner/utils"
	"urlshortner/utils/logger"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type Repository interface {
	GetLongUrl(shortUrl string) (string, error)
	GetShortUrl(shortUrl CreateShortUrl) (string, error)
	CreateShortUrl(shortUrl CreateShortUrl) (string, error)
}

type repo struct {
	DB *bun.DB
}

func NewShortUrlRepo(db *bun.DB) *repo {
	return &repo{db}
}

func (r *repo) GetShortUrl(shortUrl CreateShortUrl) (string, error) {
	defer logger.Log.Sync()

	ShortUrl := &models.ShortUrl{}
	err := r.DB.NewSelect().Model(ShortUrl).ColumnExpr("*").Where("long_url=?", shortUrl.URL).Scan(context.Background())
	if err != nil {
		return "", err
	}
	if ShortUrl.ShortURL == "" {
		logger.Log.Error("url not found", zap.Any("message", "empty url"))
		return "", utils.CustomAPIErr{Code: http.StatusNotFound, Err: errors.New("url not found")}
	}
	url := dotenv.Global.ShortURLBaseURI + "/" + ShortUrl.ShortURL

	return url, nil
}

func (r *repo) GetLongUrl(shortUrl string) (string, error) {
	defer logger.Log.Sync()

	ShortUrl := &models.ShortUrl{}

	err := r.DB.NewSelect().Model(ShortUrl).ColumnExpr("*").Where("short_url=?", shortUrl).Scan(context.Background())
	if err != nil {
		return "", err
	}
	if ShortUrl.LongURL == "" {
		logger.Log.Error("url not found", zap.Any("shortUrl", shortUrl))
		return "", utils.CustomAPIErr{Code: http.StatusNotFound, Err: errors.New("url not found")}
	}

	return ShortUrl.LongURL, nil
}

func (r *repo) CreateShortUrl(LongUrl CreateShortUrl) (string, error) {

	ShortUrl := &models.ShortUrl{LongURL: LongUrl.URL}
	_, err := r.DB.NewInsert().Model(ShortUrl).Returning("*").Exec(context.Background())

	if err != nil {
		return "", err
	}
	url := dotenv.Global.ShortURLBaseURI + "/" + ShortUrl.ShortURL
	return url, nil
}
