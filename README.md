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
* Create a `Paginator` instance with: your store, limit, offset and options
* Call the `paginator.Page()` method to get a `paging.Page` instance
* Call the `page.Previous()` method to get the previous page instance
* Call the `page.Next()` method to get the next page instance

Example with GORM:

```go
// Before anything... create a slice for your GORM models.
// Let's assume we have 100 users in our database.
users := []User{}

// Step 1: create the store. It takes your database connection pointer, a
// pointer to models and the GORM "ORDER BY" string.
store, err := paging.NewGORMStore(&db, &users)
if err != nil {
        log.Fatal(err)
}

// Step 2: create options. Here, we use the default ones (see below).
options := paging.NewOptions()

// Step 3: create a paginator instance and pass your store, your current HTTP
// request and your options as arguments.
paginator := paging.NewPaginator(store, 20, 0, options)

// You can also use the NewRequestPaginator initializer to auto-handle
// pagination for a http.Request instance.
paginator = paging.NewRequestPaginator(store, request, options)

// Step 4: calls the paginator.Page() method to get the page instance.
page, err := paginator.Page()
if err != nil {
        log.Fatal(err)
}

// Your page instance contains everything you need.
assert.True(int64(20), page.Limit)
assert.True(int64(0), page.Offset)
assert.True(int64(100), page.Count)
assert.False(page.PreviousURI.Valid) // It's a null string because no previous page
assert.True(page.NextURI.Valid)
assert.Equal( "?limit=20&offset=20", page.NextURI.String)

// And our "users" slice are now populated with 20 users order by name.
assert.Equal(20, len(users))

// Now get the next page.
nextPage, err := page.Next()
if err != nil {
        log.Fatal(err)
}

// Or the previous page.
previousPage, err := page.Previous()
if err != nil {
        log.Fatal(err)
}
```

Paginator options are:

* `DefaultLimit` (`int64`): the number of items per page (defaults to `20`)
* `LimitKeyName` (`string`): the query string key name for limit (defaults to `limit`)
* `OffsetKeyName` (`string`): the query string key name for offset (defaults to `offset`)

## Contributing

* Ping us on twitter [@thoas](https://twitter.com/thoas), [@oibafsellig](https://twitter.com/oibafsellig)
* Fork the [project](https://github.com/ulule/paging)
* Fix [bugs](https://github.com/ulule/paging/issues)
