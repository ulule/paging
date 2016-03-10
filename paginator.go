package paging

import (
	"errors"
	"net/http"

	"github.com/guregu/null"
)

// -----------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------

// Options are paginator options.
type Options struct {
	// DefaultLimit is the default number of items per page.
	DefaultLimit int64
	// MaxLimit is the maximum limit that can be set.
	MaxLimit int64
	// LimitKeyName is the query string key name for the limit.
	LimitKeyName string
	// OffsetKeyName is the query string key name for the offset.
	OffsetKeyName string
}

// NewOptions returns defaults options.
func NewOptions() *Options {
	return &Options{
		DefaultLimit:  int64(DefaultLimit),
		LimitKeyName:  DefaultLimitKeyName,
		OffsetKeyName: DefaultOffsetKeyName,
	}
}

// -----------------------------------------------------------------------------
// Page
// -----------------------------------------------------------------------------

// Page contains pagination information.
type Page struct {
	paginator *Paginator

	Limit       int64       `json:"limit"`
	Offset      int64       `json:"offset"`
	Count       int64       `json:"total_count"`
	PreviousURI null.String `json:"previous"`
	NextURI     null.String `json:"next"`
}

// Previous returns the previous page.
func (p *Page) Previous() (*Page, error) {
	if !p.paginator.HasPrevious() {
		return nil, errors.New("No previous page")
	}

	paginator := *p.paginator
	paginator.Offset = p.Offset - p.Limit

	return paginator.Page()
}

// Next returns the next page.
func (p *Page) Next() (*Page, error) {
	if !p.paginator.HasNext() {
		return nil, errors.New("No next page")
	}

	paginator := *p.paginator
	paginator.Offset = p.Offset + p.Limit

	return paginator.Page()
}

// -----------------------------------------------------------------------------
// Paginator
// -----------------------------------------------------------------------------

// Paginator is the paginator.
type Paginator struct {
	// Store is the store that contains entities to paginate.
	Store Store
	// Options are user options.
	Options *Options
	// Limit is the pagination limit.
	Limit int64
	// Offset is the pagination offset.
	Offset int64
	// Count is the total count of items to paginate.
	Count int64
}

// NewPaginator returns a new Paginator instance.
func NewPaginator(store Store, limit, offset int, options *Options) (*Paginator, error) {
	if options == nil {
		options = NewOptions()
	}

	l, o := int64(limit), int64(offset)

	if !ValidateLimitOffset(l, o) {
		return nil, errors.New("invalid limit or offset")
	}

	return &Paginator{
		Store:   store,
		Options: options,
		Limit:   l,
		Offset:  o,
	}, nil
}

// NewRequestPaginator returns a new RequestPaginator instance.
func NewRequestPaginator(store Store, request *http.Request, options *Options) (*Paginator, error) {
	if options == nil {
		options = NewOptions()
	}

	limit := GetLimitFromRequest(request, options)
	offset := GetOffsetFromRequest(request, options)

	if !ValidateLimitOffset(limit, offset) {
		return nil, errors.New("invalid limit or offset")
	}

	return &Paginator{
		Store:   store,
		Options: options,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

// HasPrevious returns true if there is a previous page.
func (p *Paginator) HasPrevious() bool {
	if (p.Offset - p.Limit) < 0 {
		return false
	}
	return true
}

// HasNext retourns true if has next page.
func (p *Paginator) HasNext() bool {
	if (p.Offset + p.Limit) >= p.Count {
		return false
	}
	return true
}

// PreviousURI returns the previous page URI.
func (p *Paginator) PreviousURI() null.String {
	if !p.HasPrevious() {
		return null.NewString("", false)
	}
	return null.StringFrom(GenerateURI(p.Limit, (p.Offset - p.Limit), p.Options))
}

// NextURI returns the next page URI.
func (p *Paginator) NextURI() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}
	return null.StringFrom(GenerateURI(p.Limit, (p.Offset + p.Limit), p.Options))
}

// Page returns the page instance and fetch items from the store.
func (p *Paginator) Page() (*Page, error) {
	if err := p.Store.Paginate(p.Limit, p.Offset, &p.Count); err != nil {
		return nil, err
	}

	return &Page{
		paginator:   p,
		Limit:       p.Limit,
		Offset:      p.Offset,
		PreviousURI: p.PreviousURI(),
		NextURI:     p.NextURI(),
		Count:       p.Count,
	}, nil
}
