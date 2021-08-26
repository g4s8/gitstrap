package gitstrap

import (
	"testing"
	"github.com/google/go-github/v38/github"
	m "github.com/g4s8/go-matchers"
)

func TestPagination(t *testing.T) {
	assert := m.Assert(t)
	pag := new(pagination)
	opts := new(github.ListOptions)
	assert.That("Moved to first page", pag.moveNext(opts), m.Is(true))
	assert.That("First move updates options", opts.Page, m.Is(1))
	rsp := &github.Response{
		NextPage: 1,
		LastPage: 2,
	}
	pag.update(rsp)
	assert.That("Moved to next page", pag.moveNext(opts), m.Is(true))
	assert.That("Next page updates options", opts.Page, m.Is(1))
	rsp.NextPage = 2
	pag.update(rsp)
	assert.That("Moved to last page", pag.moveNext(opts), m.Is(true))
	rsp.NextPage = 0
	rsp.LastPage = 0
	pag.update(rsp)
	assert.That("Last page updates options", opts.Page, m.Is(2))
	assert.That("Can't move after last page", pag.moveNext(opts), m.Is(false))
}
