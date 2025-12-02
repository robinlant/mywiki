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

type Query struct {
	Limit uint
	Field string
	Desc  bool
}

type Store interface {
	SavePage(ctx context.Context, p *Page) error
	LoadPage(ctx context.Context, title string) (*Page, bool, error)
	LoadPages(ctx context.Context, q Query) ([]*Page, error)
}
