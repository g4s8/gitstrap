package gitstrap

import (
	"github.com/google/go-github/v38/github"
)

type pagination struct {
	next, last int
}

func (p *pagination) update(rsp *github.Response) {
	p.next = rsp.NextPage
	p.last = rsp.LastPage
}

func (p *pagination) moveNext(opts *github.ListOptions) bool {
	if opts.Page == 0 && p.next == 0 {
		// initial case
		opts.Page = 1
		return true
	}
	if p.last == 0 {
		// final case
		return false
	}
	opts.Page = p.next
	return true
}
