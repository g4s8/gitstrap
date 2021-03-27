package gitstrap

import (
	"github.com/g4s8/gitstrap/internal/spec"
)

// GetRepo repository resource
func (g *Gitstrap) GetRepo(name string, owner string) (*spec.Model, error) {
	ctx, cancel := g.newContext()
	if owner == "" {
		owner = g.me
	}
	defer cancel()
	r, _, err := g.gh.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}
	model, err := spec.NewModel(spec.KindRepo)
	if err != nil {
		panic(err)
	}
	model.Metadata.FromGithubRepo(r)
	if owner != "" {
		model.Metadata.Owner = owner
	} else {
		model.Metadata.Owner = g.me
	}
	repo := new(spec.Repo)
	repo.FromGithub(r)
	model.Spec = repo
	return model, nil
}

func (g *Gitstrap) GetOrg(name string) (*spec.Model, error) {
	ctx, cancel := g.newContext()
	defer cancel()
	o, _, err := g.gh.Organizations.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	model, err := spec.NewModel(spec.KindOrg)
	if err != nil {
		panic(err)
	}
	model.Metadata.FromGithubOrg(o)
	model.Metadata.Name = name
	model.Metadata.ID = o.ID
	org := new(spec.Org)
	org.FromGithub(o)
	model.Spec = org
	return model, nil
}
