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

type User struct {
	Number int
	Name   string
}

func TestGORMStore(t *testing.T) {
	a := assert.New(t)

	conn, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db, err := gorm.Open("sqlite3", conn)
	a.Nil(err)
	db.LogMode(false)
	db.DB().SetMaxIdleConns(10)

	db.DropTableIfExists(&User{})
	db.CreateTable(&User{})

	for i := 1; i <= 100; i++ {
		db.Create(&User{
			Number: i,
			Name:   fmt.Sprintf("user-%d", i),
		})
	}

	request, _ := http.NewRequest("GET", "http://example.com", nil)

	users := []User{}

	q := db.Model(&User{})
	q = q.Order("number desc")

	store, err := NewGORMStore(q, &users)
	a.Nil(err)

	options := NewOptions()

	paginator, err := NewPaginator(store, request, options)
	a.Nil(err)

	page, err := paginator.Page()
	a.Nil(err)

	a.Equal(int64(20), page.Limit)
	a.Equal(len(users), 20)
	a.Equal(int64(0), page.Offset)
	a.Equal(int64(100), page.Count)
	a.False(page.Previous.Valid) // null
	a.Equal("?limit=20&offset=20", page.Next.String)

	// Check order desc
	a.Equal(100, users[0].Number)

	//
	// Next
	//

	request, _ = http.NewRequest("GET", page.Next.String, nil)

	paginator, err = NewPaginator(store, request, options)
	a.Nil(err)

	page, err = paginator.Page()
	a.Nil(err)

	a.Equal(int64(20), page.Limit)
	a.Equal(len(users), 20)
	a.Equal(int64(20), page.Offset)
	a.Equal(int64(100), page.Count)
	a.Equal("?limit=20&offset=0", page.Previous.String)
	a.Equal("?limit=20&offset=40", page.Next.String)

	// Check order desc
	a.Equal(80, users[0].Number)
}
