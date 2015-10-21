package paging

import (
	"fmt"
	"net/http"
	"strconv"
)

// ValidateLimitOffset returns true if limit/offset values are valid
func ValidateLimitOffset(limit int64, offset int64) bool {
	values := []int64{limit, offset}

	for _, v := range values {
		if v < 0 {
			return false
		}
	}

	return true
}

// GetLimitFromRequest returns current limit.
func GetLimitFromRequest(request *http.Request, options *Options) int64 {
	var (
		limit int64
		err   error
	)

	requestLimit := request.URL.Query().Get(options.LimitKeyName)

	if requestLimit != "" {
		limit, err = strconv.ParseInt(requestLimit, 10, 64)
		if err != nil {
			limit = options.DefaultLimit
		}
	} else {
		limit = options.DefaultLimit
	}

	return limit
}

// GetOffsetFromRequest returns current offset.
func GetOffsetFromRequest(request *http.Request, options *Options) int64 {
	var (
		offset int64
		err    error
	)

	requestOffset := request.URL.Query().Get(options.OffsetKeyName)

	if requestOffset != "" {
		offset, err = strconv.ParseInt(requestOffset, 10, 64)
		if err != nil {
			offset = 0
		}
	} else {
		offset = 0
	}

	return offset
}

// GenerateURI generates the pagination URI.
func GenerateURI(limit int64, offset int64, options *Options) string {
	return fmt.Sprintf(
		"?%s=%d&%s=%d",
		options.LimitKeyName,
		limit,
		options.OffsetKeyName,
		offset)
}
