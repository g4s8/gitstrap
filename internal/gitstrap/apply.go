package gitstrap

import (
	"errors"
	"fmt"

	gh "github.com/g4s8/gitstrap/internal/github"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/google/go-github/v33/github"
)

// Apply specification
func (g *Gitstrap) Apply(m *spec.Model) error {
	switch m.Kind {
	case spec.KindRepo:
		return g.applyRepo(m)
	case spec.KindHook:
		return g.applyHook(m)
	default:
		return fmt.Errorf("Unsupported yet %s", m.Kind)
	}
}

func (g *Gitstrap) applyRepo(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	repo := new(spec.Repo)
	if err := m.GetSpec(repo); err != nil {
		return err
	}
	meta := m.Metadata
	owner := g.getOwner(m)
	name := meta.Name
	exist, err := gh.RepoExist(g.gh, ctx, owner, name)
	if err != nil {
		return err
	}
	if !exist {
		return g.createRepo(m)
	}
	gr := new(github.Repository)
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

var errHookSelectorEmpty = errors.New("Hook selector is empty: requires repository or organization")

func (g *Gitstrap) applyHook(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	owner := g.getOwner(m)
	hook := new(spec.Hook)
	if err := m.GetSpec(hook); err != nil {
		return err
	}
	if m.Metadata.ID == nil {
		return g.createHook(m)
	}
	ghook := new(github.Hook)
	if err := hook.ToGithub(ghook); err != nil {
		return err
	}
	ghook.ID = m.Metadata.ID
	var err error
	if hook.Selector.Repository != "" {
		ghook, _, err = g.gh.Repositories.EditHook(ctx, owner, hook.Selector.Repository, *m.Metadata.ID, ghook)
	} else if hook.Selector.Organization != "" {
		ghook, _, err = g.gh.Organizations.EditHook(ctx, hook.Selector.Organization, *m.Metadata.ID, ghook)
	} else {
		err = errHookSelectorEmpty
	}
	if err != nil {
		return err
	}
	if err := hook.FromGithub(ghook); err != nil {
		return err
	}
	m.Spec = hook
	return nil
}
