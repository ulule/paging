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
	// CursorKeyName is the query string key name for the offset.
	CursorKeyName string
}

// NewOptions returns defaults options.
func NewOptions() *Options {
	return &Options{
		DefaultLimit:  int64(DefaultLimit),
		LimitKeyName:  DefaultLimitKeyName,
		OffsetKeyName: DefaultOffsetKeyName,
		CursorKeyName: DefaultCursorKeyName,
	}
}
