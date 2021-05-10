package paging

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	db *gorm.DB
)

const refDate = 1484652856

func init() {
	var err error

	conn, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open("sqlite3", conn)
	if err != nil {
		panic(err)
	}

	db.LogMode(false)
	db.DB().SetMaxIdleConns(10)
}

func rebuildDB() {
	db.DropTableIfExists(&User{})
	db.CreateTable(&User{})
	timeRef := time.Unix(0, refDate*1000000000)
	min := 100
	for i := 1; i <= 100; i++ {
		db.Create(&User{
			ID:           i,
			Number:       i,
			Name:         fmt.Sprintf("user-%d", i),
			DateCreation: timeRef.Add(time.Duration(-min) * time.Minute),
		})
		min = min - 1
	}
}

type User struct {
	ID           int
	Number       int
	Name         string
	DateCreation time.Time
}

func TestGORMStore_OffsetPaginator(t *testing.T) {
	is := assert.New(t)

	rebuildDB()

	request, _ := http.NewRequest("GET", "http://example.com", nil)

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("number desc")

	store, err := NewGORMStore(q, &users)
	is.Nil(err)

	options := NewOptions()

	paginator, err := NewOffsetPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), paginator.Offset)
	is.Equal(int64(100), paginator.Count)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&offset=20", paginator.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)

	//
	// Next
	//

	request, _ = http.NewRequest("GET", paginator.NextURI.String, nil)

	paginator, err = NewOffsetPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(20), paginator.Offset)
	is.Equal(int64(100), paginator.Count)
	is.Equal("?limit=20&offset=0", paginator.PreviousURI.String)
	is.Equal("?limit=20&offset=40", paginator.NextURI.String)

	//
	// Now, previous.
	//

	request, _ = http.NewRequest("GET", paginator.PreviousURI.String, nil)
	paginator, err = NewOffsetPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), paginator.Offset)
	is.Equal(int64(100), paginator.Count)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&offset=20", paginator.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)

	// Modify default limit
	options.SetDefaultLimit(10)
	request, _ = http.NewRequest("GET", paginator.PreviousURI.String, nil)
	paginator, err = NewOffsetPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(10), paginator.Limit)
	is.Equal(len(users), 10)
	is.Equal(int64(0), paginator.Offset)
	is.Equal(int64(100), paginator.Count)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=10&offset=10", paginator.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)
}

func TestGORMStore_CursorPaginator(t *testing.T) {
	is := assert.New(t)

	rebuildDB()

	request, _ := http.NewRequest("GET", "http://example.com?limit=20&since=0", nil)

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("number asc")

	store, err := NewGORMStore(q, &users)
	is.Nil(err)

	options := NewOptions()

	paginator, err := NewCursorPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), paginator.Cursor)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since=20", paginator.NextURI.String)

	//
	// Next with uri
	//

	request, _ = http.NewRequest("GET", paginator.NextURI.String, nil)

	paginator, err = NewCursorPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(20), paginator.Cursor)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since=40", paginator.NextURI.String)

	//
	// Next with method
	//

	np, err := paginator.Next()
	is.Nil(err)
	nextPaginator := np.(*CursorPaginator)

	is.Nil(err)

	is.Equal(int64(20), nextPaginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int(40), nextPaginator.Cursor)
	is.False(nextPaginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since=60", nextPaginator.NextURI.String)

	//
	// test previous
	//

	pp, err := nextPaginator.Previous()
	is.Nil(pp)
	is.NotNil(err)

	// Check order asc
	is.Equal(41, users[0].Number)
}

func TestGORMStore_CursorPaginator_Date(t *testing.T) {
	is := assert.New(t)

	rebuildDB()

	request, _ := http.NewRequest("GET", fmt.Sprintf("http://example.com?limit=20&since-date=%d", refDate), nil)

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("date_creation desc")

	store, err := NewGORMStore(q, &users)
	is.Nil(err)

	options := NewOptions()
	options.CursorOptions.Mode = DateModeCursor
	options.CursorOptions.DBName = "date_creation"
	options.CursorOptions.StructName = "DateCreation"
	options.CursorOptions.KeyName = "since-date"
	options.CursorOptions.Reverse = true

	paginator, err := NewCursorPaginator(store, request, options)
	is.Nil(err)
	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since-date=1484651656", paginator.NextURI.String)
	is.Equal(100, users[0].Number)

	// //
	// // Next
	// //

	request, _ = http.NewRequest("GET", paginator.NextURI.String, nil)

	paginator, err = NewCursorPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since-date=1484650456", paginator.NextURI.String)
	is.Equal(80, users[0].Number)

	// //
	// // Next again
	// //

	np, err := paginator.Next()
	is.Nil(err)
	nextPaginator := np.(*CursorPaginator)

	is.Equal(int64(20), nextPaginator.Limit)
	is.Equal(len(users), 20)
	is.False(nextPaginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&since-date=1484649256", nextPaginator.NextURI.String)
	is.Equal(60, users[0].Number)

	// //
	// // End of cursor
	// //

	// end with next
	np, err = paginator.Next()
	is.Nil(err)
	is.Equal(40, users[0].Number)

	np, err = np.Next()
	is.Nil(err)
	is.Equal(20, users[0].Number)

	np, err = np.Next()
	is.Error(err)

	// end with request
	request, _ = http.NewRequest("GET", "http://example.com?limit=20&since-date=1484646856", nil)

	paginator, err = NewCursorPaginator(store, request, options)
	is.Nil(err)

	err = paginator.Page()
	is.Nil(err)
	is.False(paginator.NextURI.Valid) // null
	is.Empty(users)

	// //
	// // test previous
	// //

	pp, err := nextPaginator.Previous()
	is.Nil(pp)
	is.NotNil(err)

}

func TestGORMStore_PaginateCursor_HasNext(t *testing.T) {
	is := assert.New(t)
	rebuildDB()

	var items []User
	s := GORMStore{db: db.Model(&User{}), items: &items}

	var hasnext bool
	is.NoError(s.PaginateCursor(99, 0, DefaultCursorDBName, false, &hasnext))
	is.Equal(99, len(items))
	is.True(hasnext)

	is.NoError(s.PaginateCursor(100, 0, DefaultCursorDBName, false, &hasnext))
	is.Equal(100, len(items))
	is.False(hasnext)
}
