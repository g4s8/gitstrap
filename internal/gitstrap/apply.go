package gitstrap

import (
	"context"
	"fmt"

	"github.com/g4s8/gitstrap/internal/github"
	"github.com/g4s8/gitstrap/internal/spec"
	gh "github.com/google/go-github/v33/github"
)

func (g *Gitstrap) Apply(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	switch m.Kind {
	case spec.KindRepo:
		return g.applyRepo(ctx, m)
	default:
		return fmt.Errorf("Unsupported yet %s", m.Kind)
	}
}

func (g *Gitstrap) applyRepo(ctx context.Context, m *spec.Model) error {
	repo := new(spec.Repo)
	if err := m.RepoSpec(repo); err != nil {
		return err
	}
	meta := m.Metadata
	owner := meta.Owner
	if owner == "" {
		owner = g.me
	}
	name := meta.Name
	exist, err := github.RepoExist(g.gh, ctx, owner, name)
	if err != nil {
		return err
	}
	if !exist {
		return g.createRepo(ctx, m)
	}
	gr := new(gh.Repository)
	if err := repo.ToGithub(gr); err != nil {
		return err
	}
	gr.ID = meta.ID
	gr, _, err = g.gh.Repositories.Edit(ctx, owner, name, gr)
	if err != nil {
		return err
	}
	repo.FromGithub(gr)
	m.Spec = repo
	m.Metadata.FromGithubRepo(gr)
	m.Metadata.Owner = owner
	return nil
}
