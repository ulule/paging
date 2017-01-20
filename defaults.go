package paging

const (
	// DefaultLimit is the default number of items per page.
	DefaultLimit = 20

	// DefaultLimitKeyName is the request key name.
	DefaultLimitKeyName = "limit"

	// DefaultOffsetKeyName is the request offset key name.
	DefaultOffsetKeyName = "offset"

	// DefaultCursorKeyName is the request cursor key name.
	DefaultCursorKeyName = "since"

	// DefaultCursorDBName is the default cursor db field name
	DefaultCursorDBName = "id"

	// DefaultCursorStructName is the default cursor struct field name
	DefaultCursorStructName = "ID"
)
