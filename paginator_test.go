package paging

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCursorPaginator(t *testing.T) {
	is := assert.New(t)

	now := time.Unix(1548002340, 0)
	since := now.UnixNano() / 1e6 // Date.now() in javascript returns a timestamp in milliseconds
	v := url.Values{"since": []string{strconv.FormatInt(since, 10)}}
	opts := NewOptions()
	opts.CursorOptions.Mode = DateModeCursor
	p, err := NewCursorPaginator(nil, &http.Request{URL: &url.URL{RawQuery: v.Encode()}}, opts)
	is.NoError(err)

	got, ok := p.Cursor.(time.Time)
	is.True(ok)
	is.True(got.After(now), "%v should be after %v", got, now)
}
