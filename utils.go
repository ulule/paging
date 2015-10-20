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

// GenerateURI generates the pagination URI.
func GenerateURI(limit int64, offset int64) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

// GetLimitFromRequest returns current limit.
func GetLimitFromRequest(request *http.Request, defaultLimit int64) int64 {

	var (
		limit int64
		err   error
	)

	requestLimit := request.URL.Query().Get(LimitKeyName)

	if requestLimit != "" {
		limit, err = strconv.ParseInt(requestLimit, 10, 64)
		if err != nil {
			limit = defaultLimit
		}
	} else {
		limit = defaultLimit
	}

	return limit
}

// GetOffsetFromRequest returns current offset.
func GetOffsetFromRequest(request *http.Request) int64 {
	var (
		offset int64
		err    error
	)

	requestOffset := request.URL.Query().Get(OffsetKeyName)

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
