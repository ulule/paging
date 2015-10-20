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

### GORM Paginator

This package includes a ready-to-go paginator for [GORM](https://github.com/jinzhu/gorm).

To use it, it's fairly simple:

```go

// For example, you want to paginate a list of users.
users := []User{}

// "db" is your database connection pointer.
// "request" is your HTTP request.
// "users" is a pointer to your user list.
// "20" is the default limit.
paginator := paging.NewGORMPaginator(&db, request, &users, 20)

// That's all. Now, you can get the current page with the Pager.Page() method.
page, err := paginator.Page()

// You can also use the Pager interface methods directly.
hasPrevious := paginator.HasPrevious()
previousURL := paginator.Previous()
hasNext := paginator.HasNext()
nextURL := pagiantor.Next()
```

## Contributing

* Ping us on twitter [@thoas](https://twitter.com/thoas), [@oibafsellig](https://twitter.com/oibafsellig)
* Fork the [project](https://github.com/ulule/paging)
* Fix [bugs](https://github.com/ulule/paging/issues)
