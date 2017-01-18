package paging

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
	// CursorMode set type of cursor, an ID or a Date (time.Time)
	CursorMode string
	// CursorKeyName is the query string key name for the offset.
	CursorKeyName string
	// CursorDBName is the cursor bdd field name
	CursorDBName string
	// CursorStructName is the cursor struct field name
	CursorStructName string
	// CursorReverse turn true to work with DESC request
	CursorReverse bool
}

// NewOptions returns defaults options.
func NewOptions() *Options {
	return &Options{
		DefaultLimit:     int64(DefaultLimit),
		LimitKeyName:     DefaultLimitKeyName,
		OffsetKeyName:    DefaultOffsetKeyName,
		CursorMode:       IDModeCursor,
		CursorKeyName:    DefaultCursorKeyName,
		CursorDBName:     DefaultCursorDBName,
		CursorStructName: DefaultCursorStructName,
		CursorReverse:    false,
	}
}
