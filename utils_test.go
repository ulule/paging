package paging

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLimitOffset(t *testing.T) {
	is := assert.New(t)

	is.True(ValidateLimitOffset(20, 10))
	is.True(ValidateLimitOffset(1, 1))
	is.True(ValidateLimitOffset(1, 1))
	is.True(ValidateLimitOffset(1, 0))
	is.True(ValidateLimitOffset(0, 1))
	is.False(ValidateLimitOffset(-1, 1))
	is.False(ValidateLimitOffset(1, -1))
}

func TestGetLimitFromRequest(t *testing.T) {
	is := assert.New(t)

	options := NewOptions()

	// We define a default limit...
	options.DefaultLimit = 40
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	is.Equal(int64(40), GetLimitFromRequest(request, options))

	// But even we define a default limit, user must be able to override it
	// by setting it in the URL.
	options.DefaultLimit = 10
	request, _ = http.NewRequest("GET", "http://example.com/?limit=33", nil)
	is.Equal(int64(33), GetLimitFromRequest(request, options))

	// If user use a different query string key, we must set the default one.
	options.DefaultLimit = 20
	options.LimitKeyName = "perpage"
	request, _ = http.NewRequest("GET", "http://example.com/?limit=33", nil)
	is.Equal(int64(20), GetLimitFromRequest(request, options))

	// Now, let's check with a good query string key.
	options.DefaultLimit = 20
	options.LimitKeyName = "perpage"
	request, _ = http.NewRequest("GET", "http://example.com/?perpage=56", nil)
	is.Equal(int64(56), GetLimitFromRequest(request, options))

	// We restrict with a max limit
	options.MaxLimit = 15
	options.LimitKeyName = "limit"
	request, _ = http.NewRequest("GET", "http://example.com/?limit=100", nil)
	is.Equal(int64(15), GetLimitFromRequest(request, options))
}

func TestGetOffsetFromRequest(t *testing.T) {
	is := assert.New(t)

	options := NewOptions()

	// No offset in URL, we should have 0.
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	is.Equal(int64(0), GetOffsetFromRequest(request, options))

	// Offset in URL, let's get it.
	request, _ = http.NewRequest("GET", "http://example.com?offset=42", nil)
	is.Equal(int64(42), GetOffsetFromRequest(request, options))

	// If user use a different query string key, we must set the default one.
	request, _ = http.NewRequest("GET", "http://example.com?offshore=90", nil)
	is.Equal(int64(0), GetOffsetFromRequest(request, options))

	// Now, let's check with a good query string key.
	options.OffsetKeyName = "mayday"
	request, _ = http.NewRequest("GET", "http://example.com?mayday=90", nil)
	is.Equal(int64(90), GetOffsetFromRequest(request, options))
}

func TestGetCursorFromRequest(t *testing.T) {
	is := assert.New(t)

	options := NewOptions()

	// No cursor in URL, we should have 0.
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	is.Equal(int64(0), GetCursorFromRequest(request, options))

	// Offset in URL, let's get it.
	request, _ = http.NewRequest("GET", "http://example.com?since=42", nil)
	is.Equal(int64(42), GetCursorFromRequest(request, options))

	// If user use a different query string key, we must set the default one.
	request, _ = http.NewRequest("GET", "http://example.com?offshore=90", nil)
	is.Equal(int64(0), GetCursorFromRequest(request, options))

	// Now, let's check with a good query string key.
	options.CursorOptions.KeyName = "mayday"
	request, _ = http.NewRequest("GET", "http://example.com?mayday=90", nil)
	is.Equal(int64(90), GetCursorFromRequest(request, options))
}

func TestGenerateOffsetURI(t *testing.T) {
	is := assert.New(t)

	options := NewOptions()
	is.Equal("?limit=10&offset=40", GenerateOffsetURI(int64(10), int64(40), options))

	options.LimitKeyName = "l"
	options.OffsetKeyName = "o"
	is.Equal("?l=14&o=60", GenerateOffsetURI(int64(14), int64(60), options))
}

func TestGenerateCursorURI(t *testing.T) {
	is := assert.New(t)

	options := NewOptions()
	is.Equal("?limit=10&since=40", GenerateCursorURI(int64(10), int64(40), options))

	options.LimitKeyName = "l"
	options.CursorOptions.KeyName = "o"
	is.Equal("?l=14&o=60", GenerateCursorURI(int64(14), int64(60), options))
}

func Test_GetLastElementField(t *testing.T) {
	last := getLastElementField(
		[]struct{ Fieldname int }{
			{Fieldname: 1},
			{Fieldname: 2},
			{Fieldname: 3}},
		"Fieldname")
	assert.New(t).Equal(3, last)
}

func Test_GetLastElementField_Int(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}
		expected := "can't get last element of a value of type int"
		if r.(string) != expected {
			t.Fatalf("expected %q, got %q", expected, r)
		}
	}()
	getLastElementField(1, "fieldname")
}

func Test_GetLastElementField_IntSlice(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}
		expected := "can't get fieldname \"fieldname\" of an element of type int"
		if r.(string) != expected {
			t.Fatalf("expected %q, got %q", expected, r)
		}
	}()
	getLastElementField([]int{1}, "fieldname")
}

func Test_GetLen(t *testing.T) {
	assert.New(t).Equal(3, getLen([]int{1, 2, 3}))
}

func Test_PopLastElement(t *testing.T) {
	array := &[]int{1, 2, 3}
	last, remaining := popLastElement(array)

	is := assert.New(t)
	is.Equal(3, last)
	is.Equal(&[]int{1, 2}, remaining)
}
