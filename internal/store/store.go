package store

import "context"

type Page struct {
	Title string
	Body  []byte
}

type Store interface {
	SavePage(ctx context.Context, p *Page) error
	LoadPage(ctx context.Context, title string) (*Page, bool, error)
}
