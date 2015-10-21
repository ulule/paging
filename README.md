# paging

**A small set of utilities to paginate your data in Go.**

## Installation

```bash
$ go get github.com/ulule/paging
```

## Usage

First, import the package:

```go
import "github.com/ulule/paging"
```

Now, you can access the package through `paging` namespace.

It works in four steps:

* Create a store (which is basically where your entities are stored)
* Create paginator options (or use default ones)
* Create a `Paginator` instance with: your store, your HTTP request and your options
* Call the `paginator.Page()` method to get a `paging.Page` instance

Example with GORM:

```go
// Before anything... create a slice for your GORM models.
// Let's assume we have 100 users in our database.
users := []User{}

// Step 1: create the store. It takes your database connection pointer, a
// pointer to models and the GORM "ORDER BY" string.
store := paging.NewGORMStore(&db, &users, "name asc")

// Step 2: create options. Here, we use the default ones (see below).
options := paging.NewOptions()

// Step 3: create a paginator instance and pass your store, your current HTTP
// request and your options as arguments.
paginator := paging.NewPaginator(store, request, options)

// Step 4: calls the paginator.Page() method to get the page instance.
page, err := paginator.Page()
if err != nil {
        log.Fatal("Oops")
}

// Your page instance contains everything you need.
assert.True(int64(20), page.Limit)
assert.True(int64(0), page.Offset)
assert.True(int64(100), page.Count)
assert.False(page.Previous.Valid) // It's a null string because no previous page
assert.True(page.Next.Valid)
assert.Equal( "?limit=20&offset=20", page.Next.String)

// And our "users" slice are now populated with 20 users order by name.
assert.Equal(20, len(users))
```

Paginator options are:

* `DefaultLimit` (`int64`): the number of items per page (defaults to `20`)
* `LimitKeyName` (`string`): the query string key name for limit (defaults to `limit`)
* `OffsetKeyName` (`string`): the query string key name for offset (defaults to `offset`)

## Contributing

* Ping us on twitter [@thoas](https://twitter.com/thoas), [@oibafsellig](https://twitter.com/oibafsellig)
* Fork the [project](https://github.com/ulule/paging)
* Fix [bugs](https://github.com/ulule/paging/issues)
