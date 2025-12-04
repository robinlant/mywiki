package web

import (
	"github.com/robinlant/mywiki/wiki/internal/quote"
	"github.com/robinlant/mywiki/wiki/internal/store"
)

//TODO rework this structs to use composition

type Display struct {
	Display  string
	ViewHref string
	Page     *store.Page
}

type BasePageData struct {
	Title string
	Quote *quote.Quote
}

type SearchPageData struct {
	BaseData BasePageData
	Displays []Display
	Page     uint
	Limit    uint
	Search   string
}

type RootPageData struct {
	BaseData  BasePageData
	Displays  []Display
	GotoHref  string
	GotoParam string
	Quote     quote.Quote
}

type ViewPageData struct {
	BaseData  BasePageData
	EditHref  string
	Page      *store.Page
	Exist     bool
	GotoHref  string
	GotoParam string
}

type EditPageData struct {
	BaseData BasePageData
	Display  string
	Page     *store.Page
	SaveHref string
	Exists   bool
}
