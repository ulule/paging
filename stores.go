package paging

import "github.com/jinzhu/gorm"

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
	order string
}

// NewGORMStore returns a new GORM store instance.
func NewGORMStore(db *gorm.DB, items interface{}, order string) (*GORMStore, error) {
	return &GORMStore{
		db:    db,
		items: items,
		order: order,
	}, nil
}

// Paginate paginates items from the store and update page instance.
func (s *GORMStore) Paginate(limit, offset int64, count *int64) error {
	q := s.db
	q = q.Limit(limit)
	q = q.Offset(offset)
	q = q.Order(s.order)
	q = q.Find(s.items)
	q = q.Limit(-1)
	q = q.Offset(-1)

	if err := q.Count(count).Error; err != nil {
		return err
	}

	return nil
}
