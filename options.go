package paging

// -----------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------

// Options are paginator options
type Options struct {
	// DefaultLimit is the default number of items per page
	DefaultLimit int64
	// MaxLimit is the maximum limit that can be set
	MaxLimit int64
	// LimitKeyName is the query string key name for the limit
	LimitKeyName string
	// OffsetKeyName is the query string key name for the offset
	OffsetKeyName string
	// CursorOptions
	CursorOptions *CursorOptions
}

// CursorOptions group all options about cursor pagination
type CursorOptions struct {
	// Mode set type of cursor, an ID or a Date (time.Time)
	Mode string
	// KeyName is the query string key name for the cursor
	KeyName string
	// DBName is the cursor's database column name
	DBName string
	// StructName is the cursor struct field name
	StructName string
	// Reverse turn true to work with DESC request
	Reverse bool
}

// NewOptions returns defaults options
func NewOptions() *Options {
	return &Options{
		DefaultLimit:  int64(DefaultLimit),
		LimitKeyName:  DefaultLimitKeyName,
		OffsetKeyName: DefaultOffsetKeyName,
		CursorOptions: &CursorOptions{
			Mode:       IDModeCursor,
			KeyName:    DefaultCursorKeyName,
			DBName:     DefaultCursorDBName,
			StructName: DefaultCursorStructName,
			Reverse:    false,
		},
	}
}
