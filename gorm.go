package paging

import (
	"errors"
	"net/http"

	"github.com/guregu/null"
	"github.com/jinzhu/gorm"
)

// GORMPaginator struct.
type GORMPaginator struct {
	DB      *gorm.DB
	Request *http.Request
	Items   interface{}
	Limit   int64
	Offset  int64
	Count   int64
}

// NewGORMPaginator returns a new paginator.Paginator instance.
func NewGORMPaginator(db *gorm.DB, request *http.Request, items interface{}, defaultLimit int64) (*GORMPaginator, error) {
	limit := GetLimitFromRequest(request, defaultLimit)
	offset := GetOffsetFromRequest(request)

	if !ValidateLimitOffset(limit, offset) {
		return nil, errors.New("invalid limit or offset")
	}

	return &GORMPaginator{
		DB:      db,
		Request: request,
		Items:   items,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

// HasPrevious returns true if there is a previous page.
func (p *GORMPaginator) HasPrevious() bool {
	if (p.Offset - p.Limit) < 0 {
		return false
	}

	return true
}

// Previous returns the previous page URI.
func (p *GORMPaginator) Previous() null.String {
	if !p.HasPrevious() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateURI(p.Limit, (p.Offset - p.Limit)))
}

// Next returns the next page URI.
func (p *GORMPaginator) Next() null.String {
	if !p.HasNext() {
		return null.NewString("", false)
	}

	return null.StringFrom(GenerateURI(p.Limit, (p.Offset + p.Limit)))
}

// HasNext retourns true if has next page.
func (p *GORMPaginator) HasNext() bool {
	if (p.Offset + p.Limit) >= p.Count {
		return false
	}

	return true
}

// Page returns a page instance.
func (p *GORMPaginator) Page() (*Page, error) {
	err := p.DB.Limit(p.Limit).Offset(p.Offset).Find(p.Items).Limit(-1).Offset(-1).Count(&p.Count).Error
	if err != nil {
		return nil, err
	}

	return &Page{
		Limit:    p.Limit,
		Offset:   p.Offset,
		Count:    p.Count,
		Previous: p.Previous(),
		Next:     p.Next(),
	}, nil
}
