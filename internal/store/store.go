package store

import (
	"context"
	"time"
)

type Page struct {
	Title     string
	Body      []byte
	UpdatedAt time.Time
}

type SearchQuery struct {
	Search string
	Page   uint
	Limit  uint
}

func (q SearchQuery) Skip() uint {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Limit
}

type OrderQuery struct {
	Limit uint
	Field string
	Desc  bool
}

type Store interface {
	SavePage(ctx context.Context, p *Page) error
	LoadPage(ctx context.Context, title string) (*Page, bool, error)
	LoadPages(ctx context.Context, q OrderQuery) ([]*Page, error)
	SearchPages(ctx context.Context, q SearchQuery) ([]*Page, error)
}
