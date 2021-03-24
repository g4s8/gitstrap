package gitstrap

import (
	"fmt"
	"math"
	"strconv"

	"github.com/g4s8/gitstrap/internal/view"
	"github.com/google/go-github/v33/github"
)

type listResult struct {
	name   string
	public bool
	fork   bool
	stars  int
	forks  int
}

func (r *listResult) isFork() (s string) {
	if r.fork {
		s = "fork"
	}
	return
}

func (r *listResult) visibility() (s string) {
	if r.public {
		s = "public"
	} else {
		s = "private"
	}
	return
}

func (r *listResult) starsStr() string {
	if r.stars < 1000 {
		return strconv.Itoa(r.stars)
	}
	val := float64(r.stars) / 1000
	if val < 10 {
		return fmt.Sprintf("%.1fK", val)
	}
	return fmt.Sprintf("%dK", int(math.Floor(val)))
}

func (r *listResult) forksStr() string {
	if r.forks < 1000 {
		return strconv.Itoa(r.stars)
	}
	val := float64(r.forks) / 1000
	if val < 10 {
		return fmt.Sprintf("%.1fK", val)
	}
	return fmt.Sprintf("%dK", int(math.Floor(val)))
}

func (r *listResult) PrintOn(p view.Printer) {
	p.Print(fmt.Sprintf("| %40s | %4s | %5s | %8s ★ | %8s ⎇ |",
		r.name, r.isFork(), r.visibility(),
		r.starsStr(), r.forksStr()))
}

// List of repositories
func (g *Gitstrap) List(filter ListFilter, owner string) (<-chan view.Printable, <-chan error) {
	if filter == nil {
		filter = LfNop
	}
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		opts := new(github.RepositoryListOptions)
	PAGINATION:
		list, rsp, err := g.gh.Repositories.List(ctx, owner, opts)
		if err != nil {
			errs <- err
			return
		}
		for _, item := range list {
			entry := new(listResult)
			entry.name = item.GetFullName()
			entry.public = !item.GetPrivate()
			entry.fork = item.GetFork()
			entry.stars = item.GetStargazersCount()
			entry.forks = item.GetForksCount()
			if filter.check(entry) {
				res <- entry
			}
		}
		if rsp.NextPage < rsp.LastPage {
			opts.Page = rsp.NextPage
			goto PAGINATION
		}

	}()
	return res, errs
}
