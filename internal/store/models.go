package store

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

type News struct {
	ID        uuid.UUID `json:"id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Source    *url.URL  `json:"source"`
	Tags      []string  `json:"tags"`
}
