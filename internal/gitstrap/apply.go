package gitstrap

import (
	"context"
	"fmt"

	"github.com/g4s8/gitstrap/internal/github"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	gh "github.com/google/go-github/v33/github"
)

func (g *Gitstrap) Apply(m *spec.Model) (<-chan view.Printable, <-chan error) {
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		var r view.Printable
		var err error
		switch m.Kind {
		case spec.KindRepo:
			repo := m.Spec.(*spec.Repo)
			r, err = g.applyRepo(ctx, repo, m.Metadata)
		}
		if err != nil {
			errs <- err
		} else {
			res <- r
		}
	}()
	return res, errs
}

type resRepoApply struct {
	*gh.Repository
}

func (r *resRepoApply) PrintOn(p view.Printer) {
	p.Print(fmt.Sprintf("Repository %s updated", r.GetFullName()))
}

func (g *Gitstrap) applyRepo(ctx context.Context, repo *spec.Repo, meta *spec.Metadata) (view.Printable, error) {
	owner := meta.Owner
	if owner == "" {
		owner = g.me
	}
	name := meta.Name
	exist, err := github.RepoExist(g.gh, ctx, owner, name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return g.createRepo(ctx, repo, meta)
	}
	gr := new(gh.Repository)
	if err := repo.ToGithub(gr); err != nil {
		return nil, err
	}
	gr.ID = meta.ID
	gr, _, err = g.gh.Repositories.Edit(ctx, owner, name, gr)
	if err != nil {
		return nil, err
	}
	return &resRepoApply{gr}, nil
}
