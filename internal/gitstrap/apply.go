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
	case spec.KindOrg:
		return g.applyOrg(m)
	case spec.KindTeam:
		return g.applyTeam(m)
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

func (g *Gitstrap) applyOrg(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	o := new(spec.Org)
	if err := m.GetSpec(o); err != nil {
		return err
	}
	meta := m.Metadata
	name := meta.Name
	exist, err := gh.OrgExist(g.gh, ctx, name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("Organization %v does not exist.", name)
	}
	org := new(github.Organization)
	if err := o.ToGithub(org); err != nil {
		return err
	}
	org.ID = meta.ID
	org, _, err = g.gh.Organizations.Edit(ctx, name, org)
	if err != nil {
		return err
	}
	o.FromGithub(org)
	m.Spec = o
	m.Metadata.FromGithubOrg(org)
	return nil
}

func (g *Gitstrap) applyTeam(m *spec.Model) error {
	_, err := getSpecifiedOwner(m)
	if err != nil {
		return err
	}
	_, err = getSpecifiedName(m)
	if err == nil {
		return g.editTeamBySlug(m)
	}
	_, err = getSpecifiedID(m)
	if err == nil {
		return g.editTeamByID(m)
	}
	return g.createTeam(m)
}

func (g *Gitstrap) editTeamByID(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	ID := m.Metadata.ID
	owner := m.Metadata.Owner
	ownerID, err := gh.GetOrgIdByName(g.gh, ctx, owner)
	if err != nil {
		return err
	}
	exist, err := gh.TeamExistByID(g.gh, ctx, ownerID, *ID)
	if err != nil {
		return err
	}
	if !exist {
		return g.createTeam(m)
	}
	t := new(spec.Team)
	if err := m.GetSpec(t); err != nil {
		return err
	}
	gTeam := new(github.NewTeam)
	if err := t.ToGithub(gTeam); err != nil {
		return err
	}
	gT, _, err := g.gh.Teams.EditTeamByID(ctx, ownerID, *ID, *gTeam, false)
	if err != nil {
		return err
	}
	if err = t.FromGithub(gT); err != nil {
		return err
	}
	m.Spec = t
	m.Metadata.FromGithubTeam(gT)
	return nil
}

func (g *Gitstrap) editTeamBySlug(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	owner := m.Metadata.Owner
	slug := m.Metadata.Name
	exist, err := gh.TeamExistBySlug(g.gh, ctx, owner, slug)
	if err != nil {
		return err
	}
	if !exist {
		return g.createTeam(m)
	}
	t := new(spec.Team)
	if err := m.GetSpec(t); err != nil {
		return err
	}
	gTeam := new(github.NewTeam)
	if err := t.ToGithub(gTeam); err != nil {
		return err
	}
	gT, _, err := g.gh.Teams.EditTeamBySlug(ctx, owner, slug, *gTeam, false)
	if err != nil {
		return err
	}
	if err = t.FromGithub(gT); err != nil {
		return err
	}
	m.Spec = t
	m.Metadata.FromGithubTeam(gT)
	return nil
}
