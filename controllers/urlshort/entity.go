package shorturl

type CreateShortUrl struct {
	URL string `json:"url" validate:"required,URL"`
}
