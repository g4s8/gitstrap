package gitstrap

import (
	"context"
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/google/go-github/v33/github"
)

func (g *Gitstrap) Delete(m *spec.Model) (<-chan view.Printable, <-chan error) {
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		var rs view.Printable
		var err error
		switch m.Kind {
		case spec.KindRepo:
			rs, err = g.deleteRepo(ctx, m.Metadata)
		case spec.KindReadme:
			rspec := m.Spec.(*spec.Readme)
			rs, err = g.deleteReadme(ctx, rspec, m.Metadata)
		default:
			errs <- &errUnsupportModelKind{m.Kind}
			return
		}
		if err != nil {
			errs <- err
		} else {
			res <- rs
		}
	}()
	return res, errs
}

type repoDeleteResult struct {
	owner string
	name  string
}

func (r *repoDeleteResult) PrintOn(p view.Printer) {
	p.Print(fmt.Sprintf("Repository %s/%s deleted successfully", r.owner, r.name))
}

func (g *Gitstrap) deleteRepo(ctx context.Context, meta *spec.Metadata) (view.Printable, error) {
	owner := meta.Owner
	if owner == "" {
		owner = g.me
	}
	_, err := g.gh.Repositories.Delete(ctx, owner, meta.Name)
	if err != nil {
		return nil, err
	}
	return &repoDeleteResult{meta.Owner, meta.Name}, nil
}

type errReadmeNotExists struct {
	owner, repo string
}

func (e *errReadmeNotExists) Error() string {
	return fmt.Sprintf("README `%s/%s` doesn't exist", e.owner, e.repo)
}

func (g *Gitstrap) deleteReadme(ctx context.Context, spec *spec.Readme, meta *spec.Metadata) (view.Printable, error) {
	owner := spec.Selector.Owner
	if owner == "" {
		owner = g.me
	}
	repo, _, err := g.gh.Repositories.Get(ctx, owner, spec.Selector.Repository)
	if err != nil {
		return nil, err
	}
	msg := "README.md removed"
	if cm, ok := meta.Annotations["commitMessage"]; ok {
		msg = cm
	}
	opts := &github.RepositoryContentFileOptions{
		Message: &msg,
	}
	getopts := new(github.RepositoryContentGetOptions)
	cnt, _, rsp, err := g.gh.Repositories.GetContents(ctx, owner, repo.GetName(), "README.md", getopts)
	if rsp.StatusCode == 404 {
		return nil, &errReadmeNotExists{owner, repo.GetName()}
	}
	if err != nil {
		return nil, err
	}
	if *cnt.Type != "file" {
		return nil, &errReadmeNotFile{*cnt.Type}
	}
	opts.SHA = cnt.SHA
	if _, _, err := g.gh.Repositories.DeleteFile(ctx, owner, repo.GetName(), "README.md", opts); err != nil {
		return nil, err
	}
	return &repoDeleteResult{owner, repo.GetName()}, nil
}
