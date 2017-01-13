package paging

import "github.com/jinzhu/gorm"

// -----------------------------------------------------------------------------
// Interfaces
// -----------------------------------------------------------------------------

// Store is a store.
type Store interface {
	PaginateOffset(limit, offset int64, count *int64) error
	PaginateCursor(limit int64, cursor int64, count *int64) error
	GetItems() []Item
}

type Item interface {
	GetID() int64
}

// -----------------------------------------------------------------------------
// GORM Store
// -----------------------------------------------------------------------------

// GORMStore is the store for GORM ORM.
type GORMStore struct {
	db    *gorm.DB
	items []Item
}

// NewGORMStore returns a new GORM store instance.
func NewGORMStore(db *gorm.DB, items []Item) (*GORMStore, error) {
	return &GORMStore{
		db:    db,
		items: items,
	}, nil
}

func (s *GORMStore) GetItems() []Item {
	return s.items
}

// PaginateOffset paginates items from the store and update page instance.
func (s *GORMStore) PaginateOffset(limit, offset int64, count *int64) error {
	q := s.db
	q = q.Limit(int(limit))
	q = q.Offset(int(offset))
	q = q.Find(s.items)
	q = q.Limit(-1)
	q = q.Offset(-1)

	if err := q.Count(count).Error; err != nil {
		return err
	}

	return nil
}

// PaginateCursor paginates items from the store and update page instance for cursor pagination system.
func (s *GORMStore) PaginateCursor(limit int64, cursor int64, count *int64) error {
	q := s.db
	q = q.Limit(int(limit))
	q = q.Where("id > ?", cursor)
	q = q.Find(s.items)
	q = q.Limit(-1)

	if err := q.Count(count).Error; err != nil {
		return err
	}

	return nil
}
