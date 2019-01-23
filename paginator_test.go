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

func TestCursorPaginator_Next_DateMode(t *testing.T) {
	is := assert.New(t)

	var u *User
	is.NoError(db.DropTableIfExists(u).Error)
	is.NoError(db.CreateTable(u).Error)
	ts := time.Unix(0, 1548252003033986000) // non zero fractional part
	is.NoError(db.Create(&User{DateCreation: ts}).Error)
	is.NoError(db.Create(&User{DateCreation: ts.Add(time.Second)}).Error)

	var users []User
	s, err := NewGORMStore(db.Model(u).Order("date_creation"), &users)
	is.NoError(err)

	v := url.Values{"limit": []string{"1"}, "since": []string{"1"}}
	opts := NewOptions()
	opts.CursorOptions.Mode = DateModeCursor
	opts.CursorOptions.DBName = "date_creation"
	opts.CursorOptions.StructName = "DateCreation"
	p, err := NewCursorPaginator(s, &http.Request{URL: &url.URL{RawQuery: v.Encode()}}, opts)
	is.NoError(err)
	is.NoError(p.Page())
	is.Len(users, 1)
	is.Equal(1, users[0].ID)

	next := p.MakeNextURI()
	is.True(next.Valid)
	// the next uri cursor is the timestamp of the last element incremented by one
	is.Contains(next.String, strconv.FormatInt(users[0].DateCreation.Unix()+1, 10))
}

func TestCursorPaginator_Next_DateMode_Reverse(t *testing.T) {
	is := assert.New(t)

	var u *User
	is.NoError(db.DropTableIfExists(u).Error)
	is.NoError(db.CreateTable(u).Error)
	ts := time.Unix(0, 1548252003033986000) // non zero fractional part
	is.NoError(db.Create(&User{DateCreation: ts}).Error)
	is.NoError(db.Create(&User{DateCreation: ts.Add(time.Second)}).Error)

	var users []User
	s, err := NewGORMStore(db.Model(u).Order("date_creation desc"), &users)
	is.NoError(err)

	since := strconv.FormatInt(time.Now().Unix(), 10)
	v := url.Values{"limit": []string{"1"}, "since": []string{since}}
	opts := NewOptions()
	opts.CursorOptions.Mode = DateModeCursor
	opts.CursorOptions.DBName = "date_creation"
	opts.CursorOptions.StructName = "DateCreation"
	opts.CursorOptions.Reverse = true
	p, err := NewCursorPaginator(s, &http.Request{URL: &url.URL{RawQuery: v.Encode()}}, opts)
	is.NoError(err)
	is.NoError(p.Page())
	is.Len(users, 1)
	is.Equal(2, users[0].ID)

	next := p.MakeNextURI()
	is.True(next.Valid)
	// the next uri cursor is the timestamp of the last element
	is.Contains(next.String, strconv.FormatInt(users[0].DateCreation.Unix(), 10))
}
