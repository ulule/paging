package paging

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// ValidateLimitOffset returns true if limit and offset values are valid
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
		if options.MaxLimit > 0 && limit > options.MaxLimit {
			limit = options.MaxLimit
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

// GetCursorFromRequest returns current cursor.
func GetCursorFromRequest(request *http.Request, options *Options) int64 {
	var (
		cursor int64
		err    error
	)

	requestCursor := request.URL.Query().Get(options.CursorOptions.KeyName)

	if requestCursor != "" {
		cursor, err = strconv.ParseInt(requestCursor, 10, 64)
		if err != nil {
			cursor = 0
		}
	} else {
		cursor = 0
	}

	return cursor
}

// GenerateOffsetURI generates the pagination URI.
func GenerateOffsetURI(limit int64, offset int64, options *Options) string {
	if options == nil {
		return ""
	}
	return fmt.Sprintf(
		"?%s=%d&%s=%d",
		options.LimitKeyName,
		limit,
		options.OffsetKeyName,
		offset)
}

// GenerateCursorURI generates the pagination URI for cursor system.
func GenerateCursorURI(limit int64, cursor interface{}, options *Options) string {
	if options == nil {
		return ""
	}
	return fmt.Sprintf(
		"?%s=%d&%s=%d",
		options.LimitKeyName,
		limit,
		options.CursorOptions.KeyName,
		cursor)
}

// GetPaginationType returns the pagination type "offeset|cursor"
// (use constant CursorType or OffsetType)
// return OffsetType by default
func GetPaginationType(request *http.Request, options *Options) string {
	if options == nil {
		options = NewOptions()
	}

	if cursor := GetCursorFromRequest(request, options); cursor > 0 {
		return CursorType
	}

	return OffsetType
}

func getLastElementField(array interface{}, fieldname string) interface{} {
	value := reflect.ValueOf(array)
	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}

	if kind != reflect.Array && kind != reflect.Slice {
		panic(fmt.Sprintf("can't get last element of a value of type %T", array))
	}

	if value.Len() == 0 {
		return nil
	}

	last := value.Index(value.Len() - 1)
	if last.Kind() != reflect.Struct {
		panic(fmt.Sprintf("can't get fieldname %q of an element of type %T", fieldname, last.Interface()))
	}

	return last.FieldByName(fieldname).Interface()
}

func getLen(array interface{}) int {
	value := reflect.ValueOf(array)
	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}

	if kind != reflect.Array && kind != reflect.Slice {
		panic(fmt.Sprintf("can't get len of a value of type %T", array))
	}

	return value.Len()
}

func popLastElement(arrayPtr interface{}) (last, remaining interface{}) {
	ptr := reflect.ValueOf(arrayPtr)
	if ptr.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("expected pointer type, got %T", arrayPtr))
	}

	array := ptr.Elem()
	if array.Kind() != reflect.Array && array.Kind() != reflect.Slice {
		panic(fmt.Sprintf("can't pop last element of a value of type %T", arrayPtr))
	}

	len := array.Len()
	if len == 0 {
		return nil, arrayPtr
	}

	last = array.Index(len - 1).Interface()

	array.Set(array.Slice(0, len-1))
	remaining = array.Addr().Interface()

	return last, remaining
}
