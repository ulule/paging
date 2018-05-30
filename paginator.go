package paging

import (
	"errors"
	"net/http"
	"time"

	"github.com/guregu/null"
)

// -----------------------------------------------------------------------------
// Paginator interface
// -----------------------------------------------------------------------------

// Paginator is a paginator interface.
type Paginator interface {
	Page() error
	Previous() (Paginator, error)
	Next() (Paginator, error)
	HasPrevious() bool
	HasNext() bool
	MakePreviousURI() null.String
	MakeNextURI() null.String
}

// paginator is the abstract paginator
type paginator struct {
	// Store is the store that contains entities to paginate.
	Store Store `json:"-"`
	// Options are user options.
	Options *Options `json:"-"`
	// Request is the HTTP request
	Request *http.Request `json:"-"`

	Limit   int64       `json:"limit"`
	NextURI null.String `json:"next"`
}

// -----------------------------------------------------------------------------
// Paginator with cursor
// -----------------------------------------------------------------------------

// CursorPaginator is the paginator with cursor pagination system.
type CursorPaginator struct {
	*paginator
	Cursor      interface{} `json:"-"`
	PreviousURI null.String `json:"-"`
}

// NewCursorPaginator returns a new CursorPaginator instance.
func NewCursorPaginator(store Store, request *http.Request, options *Options) (*CursorPaginator, error) {
	if options == nil {
		options = NewOptions()
	}

	paginator := &CursorPaginator{
		&paginator{
			Store:   store,
			Options: options,
			Request: request,
			Limit:   GetLimitFromRequest(request, options),
		},
		GetCursorFromRequest(request, options),
		null.NewString("", false),
	}

	if options.CursorOptions.Mode == DateModeCursor {
		// time in cursor is standard timestamp (second)
		paginator.Cursor = time.Unix(0, GetCursorFromRequest(request, options)*1000000000)
	}

	return paginator, nil
}

// Page searches and returns the items
func (p *CursorPaginator) Page() error {
	err := p.Store.PaginateCursor(
		p.Limit,
		p.Cursor,
		p.Options.CursorOptions.DBName,
		p.Options.CursorOptions.Reverse)
	if err != nil {
		return err
	}

	p.PreviousURI = p.MakePreviousURI()
	p.NextURI = p.MakeNextURI()

	return nil
}

// Previous is not available on cursor system
func (p *CursorPaginator) Previous() (Paginator, error) {
	return nil, errors.New("No previous page")
}

// Next returns next items
func (p *CursorPaginator) Next() (Paginator, error) {
	if !p.HasNext() {
		return nil, errors.New("No next page")
	}

	np := *p
	np.Cursor = Last(p.Store.GetItems(), np.Options.CursorOptions.StructName)
	err := np.Store.PaginateCursor(
		np.Limit,
		np.Cursor,
		np.Options.CursorOptions.DBName,
		np.Options.CursorOptions.Reverse)
	if err != nil {
		return nil, err
	}

	np.NextURI = p.MakeNextURI()

	return &np, nil
}

// HasPrevious returns false, previous page is not available on cursor system
func (CursorPaginator) HasPrevious() bool {
	return false
}

// HasNext returns true if has next page.
func (CursorPaginator) HasNext() bool {
	return true
}

// MakePreviousURI returns an empty URI.
func (CursorPaginator) MakePreviousURI() null.String {
	return null.NewString("", false)
}

// MakeNextURI returns the next page URI.
func (p *CursorPaginator) MakeNextURI() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}

	nextCursor := Last(p.Store.GetItems(), p.Options.CursorOptions.StructName)

	if nextCursor == nil {
		return null.NewString("", false)
	}

	// convert to timestamp if cusror mode is Date
	if p.Options.CursorOptions.Mode == DateModeCursor {
		nextCursor = nextCursor.(time.Time).Unix()
	}

	return null.StringFrom(GenerateCursorURI(p.Limit, nextCursor, p.Options))
}

// -----------------------------------------------------------------------------
// Paginator with offset
// -----------------------------------------------------------------------------

// ErrInvalidLimitOrOffset is returned by the OffsetPaginator's Page method to
// indicate that the limit or the offset is invalid
var ErrInvalidLimitOrOffset = errors.New("invalid limit or offset")

// OffsetPaginator is the paginator with offset pagination system.
type OffsetPaginator struct {
	*paginator
	Offset      int64       `json:"offset"`
	Count       int64       `json:"total_count"`
	PreviousURI null.String `json:"previous"`
}

// NewOffsetPaginator returns a new OffsetPaginator instance.
func NewOffsetPaginator(store Store, request *http.Request, options *Options) (*OffsetPaginator, error) {
	if options == nil {
		options = NewOptions()
	}

	return &OffsetPaginator{
		&paginator{
			Store:   store,
			Options: options,
			Request: request,
			Limit:   GetLimitFromRequest(request, options),
		},
		GetOffsetFromRequest(request, options),
		0,
		null.NewString("", false),
	}, nil
}

// Page searches and returns the items
func (p *OffsetPaginator) Page() error {
	if !ValidateLimitOffset(p.Limit, p.Offset) {
		return ErrInvalidLimitOrOffset
	}

	if err := p.Store.PaginateOffset(p.Limit, p.Offset, &p.Count); err != nil {
		return err
	}

	p.PreviousURI = p.MakePreviousURI()
	p.NextURI = p.MakeNextURI()

	return nil
}

// Previous returns previous items
func (p *OffsetPaginator) Previous() (Paginator, error) {
	if !p.HasPrevious() {
		return nil, errors.New("No previous page")
	}

	paginator := *p

	paginator.Offset = p.Offset - p.Limit

	if err := paginator.Store.PaginateOffset(paginator.Limit, paginator.Offset, &paginator.Count); err != nil {
		return nil, err
	}

	paginator.PreviousURI = p.MakePreviousURI()
	paginator.NextURI = p.MakeNextURI()

	return &paginator, nil
}

// Next returns next items
func (p *OffsetPaginator) Next() (Paginator, error) {
	if !p.HasNext() {
		return nil, errors.New("No next page")
	}

	paginator := *p

	paginator.Offset = p.Offset + p.Limit

	if err := paginator.Store.PaginateOffset(paginator.Limit, paginator.Offset, &paginator.Count); err != nil {
		return nil, err
	}

	paginator.PreviousURI = p.MakePreviousURI()
	paginator.NextURI = p.MakeNextURI()

	return &paginator, nil
}

// HasPrevious returns true if there is a previous page.
func (p *OffsetPaginator) HasPrevious() bool {
	if (p.Offset - p.Limit) < 0 {
		return false
	}
	return true
}

// HasNext returns true if has next page.
func (p *OffsetPaginator) HasNext() bool {
	if (p.Offset + p.Limit) >= p.Count {
		return false
	}
	return true
}

// MakePreviousURI returns the previous page URI.
func (p *OffsetPaginator) MakePreviousURI() null.String {
	if !p.HasPrevious() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateOffsetURI(p.Limit, (p.Offset - p.Limit), p.Options))
}

// MakeNextURI returns the next page URI.
func (p *OffsetPaginator) MakeNextURI() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateOffsetURI(p.Limit, (p.Offset + p.Limit), p.Options))
}
