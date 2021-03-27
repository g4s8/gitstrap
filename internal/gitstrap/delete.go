package gitstrap

import (
	"context"
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/google/go-github/v33/github"
)

func (g *Gitstrap) Delete(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	switch m.Kind {
	case spec.KindRepo:
		return g.deleteRepo(ctx, m)
	case spec.KindReadme:
		return g.deleteReadme(ctx, m)
	default:
		return &errUnsupportModelKind{m.Kind}
	}
}

func (g *Gitstrap) deleteRepo(ctx context.Context, m *spec.Model) error {
	meta := m.Metadata
	owner := meta.Owner
	if owner == "" {
		owner = g.me
	}
	if _, err := g.gh.Repositories.Delete(ctx, owner, meta.Name); err != nil {
		return err
	}
	return nil
}

type errReadmeNotExists struct {
	owner, repo string
}

func (e *errReadmeNotExists) Error() string {
	return fmt.Sprintf("README `%s/%s` doesn't exist", e.owner, e.repo)
}

func (g *Gitstrap) deleteReadme(ctx context.Context, m *spec.Model) error {
	spec := new(spec.Readme)
	meta := m.Metadata
	if err := m.ReadmeSpec(spec); err != nil {
		return err
	}
	owner := spec.Selector.Owner
	if owner == "" {
		owner = g.me
	}
	repo, _, err := g.gh.Repositories.Get(ctx, owner, spec.Selector.Repository)
	if err != nil {
		return err
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
		return &errReadmeNotExists{owner, repo.GetName()}
	}
	if err != nil {
		return err
	}
	if *cnt.Type != "file" {
		return &errReadmeNotFile{*cnt.Type}
	}
	opts.SHA = cnt.SHA
	if _, _, err := g.gh.Repositories.DeleteFile(ctx, owner, repo.GetName(), "README.md", opts); err != nil {
		return err
	}
	return nil
}
