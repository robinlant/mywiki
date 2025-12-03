package web

import "github.com/robinlant/mywiki/internal/store"

type Display struct {
	Display  string
	ViewHref string
	Page     *store.Page
}

type RootPageData struct {
	Title    string
	Displays []Display
}

type ViewPageData struct {
	Title    string
	EditHref string
	Page     *store.Page
	Exist    bool
}

type EditPageData struct {
	Title    string
	Display  string
	Page     *store.Page
	SaveHref string
}
