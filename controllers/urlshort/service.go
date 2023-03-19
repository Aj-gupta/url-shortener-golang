package shorturl

type Service interface {
	GetLongUrl(shortUrl string) (string, error)
	GetShortUrl(shortUrl CreateShortUrl) (string, error)
	CreateShortUrl(shortUrl CreateShortUrl) (string, error)
}

type svc struct {
	repo Repository
}

func NewShortUrlSvc(repo Repository) *svc {
	return &svc{repo: repo}
}

func (s *svc) GetShortUrl(shortUrl CreateShortUrl) (string, error) {
	return s.repo.GetShortUrl(shortUrl)
}

func (s *svc) GetLongUrl(shortUrl string) (string, error) {
	return s.repo.GetLongUrl(shortUrl)
}
func (s *svc) CreateShortUrl(shortUrl CreateShortUrl) (string, error) {
	return s.repo.CreateShortUrl(shortUrl)
}
