package gitstrap

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/google/go-github/v38/github"
)

type RepoInfo struct {
	name   string
	public bool
	fork   bool
	stars  int
	forks  int
}

func (r *RepoInfo) isFork() (s string) {
	if r.fork {
		s = "fork"
	}
	return
}

func (r *RepoInfo) visibility() (s string) {
	if r.public {
		s = "public"
	} else {
		s = "private"
	}
	return
}

func (r *RepoInfo) starsStr() string {
	if r.stars < 1000 {
		return strconv.Itoa(r.stars)
	}
	val := float64(r.stars) / 1000
	if val < 10 {
		return fmt.Sprintf("%.1fK", val)
	}
	return fmt.Sprintf("%dK", int(math.Floor(val)))
}

func (r *RepoInfo) forksStr() string {
	if r.forks < 1000 {
		return strconv.Itoa(r.stars)
	}
	val := float64(r.forks) / 1000
	if val < 10 {
		return fmt.Sprintf("%.1fK", val)
	}
	return fmt.Sprintf("%dK", int(math.Floor(val)))
}

func (r *RepoInfo) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "| %40s | %4s | %7s | %8s ★ | %8s ⎇ |",
		r.name, r.isFork(), r.visibility(),
		r.starsStr(), r.forksStr())
	return int64(n), err
}

// ListRepos lists repositories
func (g *Gitstrap) ListRepos(filter ListFilter, owner string, errs chan<- error) (<-chan *RepoInfo) {
	if filter == nil {
		filter = LfNop
	}
	const (
		pageSize = 10
	)
	res := make(chan *RepoInfo, pageSize)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer cancel()
		opts := &github.RepositoryListOptions{
			Visibility: "all",
		}
		pag := new(pagination)
		opts.PerPage = pageSize
		for pag.moveNext(&opts.ListOptions) {
			list, rsp, err := g.gh.Repositories.List(ctx, owner, opts)
			if err != nil {
				errs <- err
				return
			}
			pag.update(rsp)
			for _, item := range list {
				entry := new(RepoInfo)
				entry.name = item.GetFullName()
				entry.public = !item.GetPrivate()
				entry.fork = item.GetFork()
				entry.stars = item.GetStargazersCount()
				entry.forks = item.GetForksCount()
				if filter.check(entry) {
					res <- entry
				}
			}
		}

	}()
	return res
}
