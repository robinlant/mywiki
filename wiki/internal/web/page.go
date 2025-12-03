package web

import "github.com/robinlant/mywiki/wiki/internal/store"

type Display struct {
	Display  string
	ViewHref string
	Page     *store.Page
}

type SearchPageData struct {
	Title    string
	Displays []Display
	Page     uint
	Limit    uint
	Search   string
}

type RootPageData struct {
	Title     string
	Displays  []Display
	GotoHref  string
	GotoParam string
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
	Exists   bool
}
