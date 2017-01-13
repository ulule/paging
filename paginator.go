package paging

import (
	"errors"
	"net/http"

	"github.com/guregu/null"
)

// -----------------------------------------------------------------------------
// Paginator interface
// -----------------------------------------------------------------------------

// Paginator is a paginator interface.
type Paginator interface {
	Page() (interface{}, error)
	GetItems() interface{}
	Previous() (interface{}, error)
	Next() (interface{}, error)
	HasPrevious() bool
	HasNext() bool
	MakePreviousURI() null.String
	MakeNextURI() null.String
}

// AbstractPaginator is the abstract paginator
type AbstractPaginator struct {
	// Store is the store that contains entities to paginate.
	Store Store `json:"-"`
	// Options are user options.
	Options *Options `json:"-"`
	// Request is the HTTP request
	Request *http.Request `json:"-"`

	Limit       int64       `json:"limit"`
	Count       int64       `json:"total_count"`
	PreviousURI null.String `json:"previous"`
	NextURI     null.String `json:"next"`
}

// GetItems returns the result set
func (ap *AbstractPaginator) GetItems() interface{} {
	return ap.Store.GetItems()
}

// -----------------------------------------------------------------------------
// Paginator with cursor
// -----------------------------------------------------------------------------

// PaginatorCursor is the paginator with cursor pagination system.
type PaginatorCursor struct {
	*AbstractPaginator
	Cursor int64 `json:"cursor"`
}

// NewPaginatorCursor returns a new PaginatorCursor instance.
func NewPaginatorCursor(store Store, request *http.Request, options *Options) (*PaginatorCursor, error) {
	if options == nil {
		options = NewOptions()
	}

	return &PaginatorCursor{
		&AbstractPaginator{
			Store:   store,
			Options: options,
			Request: request,
			Limit:   GetLimitFromRequest(request, options),
		},
		GetCursorFromRequest(request, options),
	}, nil
}

// Page searches and returns the items
func (p *PaginatorCursor) Page() (interface{}, error) {

	if !ValidateLimitMarker(p.Limit, p.Cursor) {
		return nil, errors.New("invalid limit or cursor")
	}

	if err := p.Store.PaginateCursor(p.Limit, p.Cursor, &p.Count); err != nil {
		return nil, err
	}

	p.PreviousURI = p.MakePreviousURI()
	p.NextURI = p.MakeNextURI()

	return p.GetItems(), nil
}

// Previous is not available on cursor system
func (p *PaginatorCursor) Previous() (interface{}, error) {
	return nil, errors.New("No previous page")
}

// Next returns next items
func (p *PaginatorCursor) Next() (interface{}, error) {
	if !p.HasNext() {
		return nil, errors.New("No next page")
	}

	// TODO get last of current list or make request
	return p.GetItems(), nil
}

// HasPrevious returns false, previous page is not available on cursor system
func (p *PaginatorCursor) HasPrevious() bool {
	return false
}

// HasNext returns true if has next page.
func (p *PaginatorCursor) HasNext() bool {
	// TODO
	return true
}

// MakePreviousURI returns an empty URI.
func (p *PaginatorCursor) MakePreviousURI() null.String {
	return null.NewString("", false)
}

// MakeNextURI returns the next page URI.
func (p *PaginatorCursor) MakeNextURI() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}

	// TODO get last of current list or make request
	return null.StringFrom(GenerateCursorURI(p.Limit, p.Cursor, p.Options))
}

// -----------------------------------------------------------------------------
// Paginator with offset
// -----------------------------------------------------------------------------

// PaginatorOffset is the paginator with offset pagination system.
type PaginatorOffset struct {
	*AbstractPaginator
	Offset int64 `json:"offset"`
}

// NewPaginatorOffset returns a new PaginatorOffset instance.
func NewPaginatorOffset(store Store, request *http.Request, options *Options) (*PaginatorOffset, error) {
	if options == nil {
		options = NewOptions()
	}

	return &PaginatorOffset{
		&AbstractPaginator{
			Store:   store,
			Options: options,
			Request: request,
			Limit:   GetLimitFromRequest(request, options),
		},
		GetOffsetFromRequest(request, options),
	}, nil
}

// Page searches and returns the items
func (p *PaginatorOffset) Page() (interface{}, error) {
	if !ValidateLimitMarker(p.Limit, p.Offset) {
		return nil, errors.New("invalid limit or offset")
	}

	if err := p.Store.PaginateOffset(p.Limit, p.Offset, &p.Count); err != nil {
		return nil, err
	}

	p.PreviousURI = p.MakePreviousURI()
	p.NextURI = p.MakeNextURI()

	return p.GetItems(), nil
}

// Previous returns previous items
func (p *PaginatorOffset) Previous() (interface{}, error) {

	if !p.HasPrevious() {
		return nil, errors.New("No previous page")
	}

	p.Offset = p.Offset - p.Limit

	p.Store.PaginateOffset(p.Limit, p.Offset, &p.Count)

	return p.GetItems(), nil
}

// Next returns next items
func (p *PaginatorOffset) Next() (interface{}, error) {
	if !p.HasNext() {
		return nil, errors.New("No next page")
	}

	p.Offset = p.Offset + p.Limit

	p.Store.PaginateOffset(p.Limit, p.Offset, &p.Count)

	return p.GetItems(), nil
}

// HasPrevious returns true if there is a previous page.
func (p *PaginatorOffset) HasPrevious() bool {
	if (p.Offset - p.Limit) < 0 {
		return false
	}
	return true
}

// HasNext returns true if has next page.
func (p *PaginatorOffset) HasNext() bool {
	if (p.Offset + p.Limit) >= p.Count {
		return false
	}
	return true
}

// MakePreviousURI returns the previous page URI.
func (p *PaginatorOffset) MakePreviousURI() null.String {
	if !p.HasPrevious() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateOffsetURI(p.Limit, (p.Offset - p.Limit), p.Options))
}

// MakeNextURI returns the next page URI.
func (p *PaginatorOffset) MakeNextURI() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateOffsetURI(p.Limit, (p.Offset + p.Limit), p.Options))
}
