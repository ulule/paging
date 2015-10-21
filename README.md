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

It works in three steps:

* Create a store (which is basically where your entities are stored);
* Create a `Paginator` instance with: your store, your HTTP request and your limit;
* Call the `paginator.Page()` method to get a `paging.Page` instance

Example with GORM:

```go
// Before anything... create a slice for your GORM models.
users := []User{}

// Step 1: create the store. It takes your database connection pointer and a
// pointer to the slices that will contains your instances.
store := paging.NewGORMStore(&db, &users)

// Step 2: create a paginator instance and pass your store, your current HTTP
// request and your limit (number of items per page) as arguments.
paginator := paging.NewPaginator(store, request, 20)

// Step 3: calls the paginator.Page() method to get the page instance.
page, err := paginator.Page()
if err != nil {
        log.Fatal("Oops")
}
```

## Contributing

* Ping us on twitter [@thoas](https://twitter.com/thoas), [@oibafsellig](https://twitter.com/oibafsellig)
* Fork the [project](https://github.com/ulule/paging)
* Fix [bugs](https://github.com/ulule/paging/issues)
