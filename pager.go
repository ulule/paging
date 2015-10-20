package paging

import "github.com/guregu/null"

// Pager is the paginator interface.
type Pager interface {
	Previous() null.String
	HasPrevious() bool
	Next() null.String
	HasNext() bool
	Page() (*Page, error)
}

// Page contains pagination information.
type Page struct {
	Limit    int64       `json:"limit"`
	Offset   int64       `json:"offset"`
	Count    int64       `json:"total_count"`
	Previous null.String `json:"previous"`
	Next     null.String `json:"next"`
}
