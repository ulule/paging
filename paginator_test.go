package paging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursorPaginator_Next(t *testing.T) {
	is := assert.New(t)
	p := CursorPaginator{count: 10, paginator: &paginator{Limit: 10}}
	is.False(p.HasNext())
	p.count = 11
	is.True(p.HasNext())
}
