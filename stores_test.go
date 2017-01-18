package paging

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	db *gorm.DB
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

	paginator, err := NewPaginatorOffset(store, request, options)
	is.Nil(err)

	_, err = paginator.Page()
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

	paginator, err = NewPaginatorOffset(store, request, options)
	is.Nil(err)

	_, err = paginator.Page()
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
	paginator, err = NewPaginatorOffset(store, request, options)
	is.Nil(err)

	_, err = paginator.Page()
	is.Nil(err)

	is.Equal(int64(20), paginator.Limit)
	is.Equal(len(users), 20)
	is.Equal(int64(0), paginator.Offset)
	is.Equal(int64(100), paginator.Count)
	is.False(paginator.PreviousURI.Valid) // null
	is.Equal("?limit=20&offset=20", paginator.NextURI.String)

	// Check order desc
	is.Equal(100, users[0].Number)
}
