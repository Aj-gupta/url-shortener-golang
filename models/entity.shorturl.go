package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ShortUrl struct {
	bun.BaseModel `bun:"table:short_urls"`
	ID            uint64    `json:"id" bun:"id,scanonly,pk"`
	ShortURL      string    `json:"short_url" bun:"short_url,scanonly"`
	LongURL       string    `json:"long_url" bun:"long_url,notnull"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,scanonly"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,scanonly"`
	ExpiresAt     time.Time `json:"expires_at" bun:"expires_at,scanonly"`
}
