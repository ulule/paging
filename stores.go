package paging

import "github.com/ulule/gorm"

// -----------------------------------------------------------------------------
// Interfaces
// -----------------------------------------------------------------------------

// Store is a store.
type Store interface {
	Paginate(limit, offset int64, count *int64) error
}

// -----------------------------------------------------------------------------
// GORM Store
// -----------------------------------------------------------------------------

// GORMStore is the store for GORM ORM.
type GORMStore struct {
	db    *gorm.DB
	items interface{}
}

// NewGORMStore returns a new GORM store instance.
func NewGORMStore(db *gorm.DB, items interface{}) (*GORMStore, error) {
	return &GORMStore{
		db:    db,
		items: items,
	}, nil
}

// Paginate paginates items from the store and update page instance.
func (s *GORMStore) Paginate(limit, offset int64, count *int64) error {
	q := s.db
	q = q.Limit(limit)
	q = q.Offset(offset)
	q = q.Find(s.items)
	q = q.Limit(-1)
	q = q.Offset(-1)

	if err := q.Count(count).Error; err != nil {
		return err
	}

	return nil
}
