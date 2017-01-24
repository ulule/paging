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

There are two types of paginators:
- `OffsetPaginator`: uses a database offset. Returns the total of elements.
- `CursorPaginator`: uses a database condition like `ID > ?` or `creation_date < ?`. Does not return the total number of items but increases performances.

It works in four steps:

* Create a store (which is basically where your entities are stored)
* Create paginator options (or use default ones)
* Create an `OffsetPaginator` or a `CursorPaginator` instance with: your store, the HTTP request, and options
* Call the `paginator.Page()` method to process the pagination
* Call the `paginator.Previous()` method to get the previous paginator instance. (Previous page isn't available for cursor pagination system)
* Call the `paginator.Next()` method to get the next paginator instance

Example with OffsetPaginator and GORM:

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
request, _ := http.NewRequest("GET", "http://example.com?limit=20&offset=0", nil)
paginator := paging.NewOffsetPaginator(store, request, options)

// Step 4: call the paginator.Page() method to get the page instance.
err := paginator.Page()
if err != nil {
        log.Fatal(err)
}

// Your paginator instance contains everything you need.
assert.True(int64(20), paginator.Limit)
assert.True(int64(0), paginator.Offset)
assert.True(int64(100), paginator.Count)
assert.False(paginator.PreviousURI.Valid) // It's a null string because no previous page
assert.True(paginator.NextURI.Valid)
assert.Equal( "?limit=20&offset=20", paginator.NextURI.String)

// And our "users" slice is now populated with 20 users ordered by name.
assert.Equal(20, len(users))

// Now get the next page.
nextPaginator, err := paginator.Next()
if err != nil {
        log.Fatal(err)
}

// Or the previous page.
previousPaginator, err := paginator.Previous()
if err != nil {
        log.Fatal(err)
}
```

Paginator options are:

* `DefaultLimit` (`int64`): the number of items per page (defaults to `20`)
* `MaxLimit` (`int64`): the maximum limit that can be set (defaults to `20`)
* `LimitKeyName` (`string`): the query string key name for limit (defaults to `limit`)
* `OffsetKeyName` (`string`): the query string key name for offset (defaults to `offset`)
* `CursorOptions.Mode` (`string`): set type of cursor, an `idCursor` or a `dateCursor` (time.Time) (defaults to `idCursor`)
* `CursorOptions.KeyName` (`string`): the query string key name for the cursor (defaults to `since`)
* `CursorOptions.DBName` (`string`): the cursor's database column name (defaults to `id`)
* `CursorOptions.StructName` (`string`): the cursor struct field name (defaults to `ID`)
* `CursorOptions.Reverse` (`bool`): if true, order is reversed (DESC) (defaults to `false`)

## Contributing

* Ping us on twitter [@thoas](https://twitter.com/thoas), [@oibafsellig](https://twitter.com/oibafsellig), [@NotDrana](https://twitter.com/notdrana)
* Fork the [project](https://github.com/ulule/paging)
* Fix [bugs](https://github.com/ulule/paging/issues)
