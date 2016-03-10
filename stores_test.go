package paging

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/gorm"
)

var (
	db gorm.DB
)

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
	for i := 1; i <= 100; i++ {
		db.Create(&User{
			Number: i,
			Name:   fmt.Sprintf("user-%d", i),
		})
	}
}

type User struct {
	Number int
	Name   string
}

func TestGORMStore_Paginator(t *testing.T) {
	is := assert.New(t)

	rebuildDB()

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("number desc")

	store, err := NewGORMStore(q, &users)
	is.Nil(err)

	options := NewOptions()

	paginator, err := NewPaginator(store, 20, 0, options)
	is.Nil(err)

	page, err := paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), page.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), page.Offset)
	is.Equal(int64(100), page.Count)
	is.False(page.PreviousURI.Valid) // null
	is.Equal("?limit=20&offset=20", page.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)

	//
	// Next
	//

	nextPage, err := page.Next()
	is.Nil(err)

	is.Equal(int64(20), nextPage.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(20), nextPage.Offset)
	is.Equal(int64(100), nextPage.Count)
	is.Equal("?limit=20&offset=0", nextPage.PreviousURI.String)
	is.Equal("?limit=20&offset=40", nextPage.NextURI.String)

	// Check order desc
	is.Equal(80, users[0].Number)

	//
	// Next again
	//

	nextPage, err = nextPage.Next()
	is.Nil(err)

	is.Equal(int64(20), nextPage.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(40), nextPage.Offset)
	is.Equal(int64(100), nextPage.Count)
	is.Equal("?limit=20&offset=20", nextPage.PreviousURI.String)
	is.Equal("?limit=20&offset=60", nextPage.NextURI.String)

	// Check order desc
	is.Equal(60, users[0].Number)
}

func TestGORMStore_RequestPaginator(t *testing.T) {
	is := assert.New(t)

	rebuildDB()

	request, _ := http.NewRequest("GET", "http://example.com", nil)

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("number desc")

	store, err := NewGORMStore(q, &users)
	is.Nil(err)

	options := NewOptions()

	paginator, err := NewRequestPaginator(store, request, options)
	is.Nil(err)

	page, err := paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), page.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), page.Offset)
	is.Equal(int64(100), page.Count)
	is.False(page.PreviousURI.Valid) // null
	is.Equal("?limit=20&offset=20", page.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)

	//
	// Next
	//

	request, _ = http.NewRequest("GET", page.NextURI.String, nil)

	paginator, err = NewRequestPaginator(store, request, options)
	is.Nil(err)

	page, err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), page.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(20), page.Offset)
	is.Equal(int64(100), page.Count)
	is.Equal("?limit=20&offset=0", page.PreviousURI.String)
	is.Equal("?limit=20&offset=40", page.NextURI.String)

	// Check order desc
	is.Equal(80, users[0].Number)
}
