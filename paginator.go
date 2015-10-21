package paging

import (
	"net/http"

	"github.com/guregu/null"
)

// -----------------------------------------------------------------------------
// Page
// -----------------------------------------------------------------------------

// Page contains pagination information.
type Page struct {
	Limit    int64       `json:"limit"`
	Offset   int64       `json:"offset"`
	Count    int64       `json:"total_count"`
	Previous null.String `json:"previous"`
	Next     null.String `json:"next"`
}

// -----------------------------------------------------------------------------
// Paginator
// -----------------------------------------------------------------------------

// Paginator is the paginator.
type Paginator struct {
	// Store is the store that contains entities to paginate.
	Store Store

	// Request is the HTTP request.
	Request *http.Request

	// Limit is the pagination limit.
	Limit int64

	// Offset is the pagination offset.
	Offset int64

	// Count is the total count of items to paginate.
	Count int64
}

// HasPrevious returns true if there is a previous page.
func (p Paginator) HasPrevious() bool {
	if (p.Offset - p.Limit) < 0 {
		return false
	}
	return true
}

// Previous returns the previous page URI.
func (p Paginator) Previous() null.String {
	if !p.HasPrevious() {
		return null.NewString("", false)
	}
	return null.StringFrom(GenerateURI(p.Limit, (p.Offset - p.Limit)))
}

// Next returns the next page URI.
func (p *Paginator) Next() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}
	return null.StringFrom(GenerateURI(p.Limit, (p.Offset + p.Limit)))
}

// HasNext retourns true if has next page.
func (p Paginator) HasNext() bool {
	if (p.Offset + p.Limit) >= p.Count {
		return false
	}
	return true
}

// Page returns the page instance and fetch items from the store.
func (p *Paginator) Page() (*Page, error) {
	if err := p.Store.Paginate(p.Limit, p.Offset, p.Count); err != nil {
		return nil, err
	}

	return &Page{
		Limit:    p.Limit,
		Offset:   p.Offset,
		Previous: p.Previous(),
		Next:     p.Next(),
	}, nil
}
