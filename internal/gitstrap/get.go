package gitstrap

import (
	"fmt"
	"strconv"

	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/google/go-github/v36/github"
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

func (g *Gitstrap) GetHooks(owner, name string) (<-chan *spec.Model, <-chan error) {
	const (
		hooksPageSize = 10
	)
	ctx, cancel := g.newContext()
	out := make(chan *spec.Model, hooksPageSize)
	errs := make(chan error)
	if owner == "" {
		owner = g.me
	}
	go func() {
		defer close(out)
		defer close(errs)
		defer cancel()

		opts := &github.ListOptions{PerPage: hooksPageSize}
		pag := new(pagination)
		for pag.moveNext(opts) {
			var (
				ghooks []*github.Hook
				rsp    *github.Response
				err    error
			)
			if name != "" {
				ghooks, rsp, err = g.gh.Repositories.ListHooks(ctx, owner, name, opts)
			} else {
				ghooks, rsp, err = g.gh.Organizations.ListHooks(ctx, owner, opts)
			}
			pag.update(rsp)
			if err != nil {
				errs <- err
				return
			}
			for _, gh := range ghooks {
				s, err := spec.NewModel(spec.KindHook)
				if err != nil {
					panic(err)
				}
				s.Metadata.Owner = owner
				s.Metadata.ID = gh.ID
				hook := new(spec.Hook)
				if name != "" {
					hook.Selector.Repository = name
				} else {
					hook.Selector.Organization = owner
				}
				if err := hook.FromGithub(gh); err != nil {
					errs <- err
					return
				}
				s.Spec = hook
				out <- s
			}
		}

	}()
	return out, errs
}

// GetTeams fetches organization teams definitions into spec channel
func (g *Gitstrap) GetTeams(org string) (<-chan *spec.Model, <-chan error) {
	const (
		perPage = 10
	)
	res := make(chan *spec.Model, perPage)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		opts := &github.ListOptions{PerPage: perPage}
		pag := new(pagination)
		for pag.moveNext(opts) {
			batch, rsp, err := g.gh.Teams.ListTeams(ctx, org, opts)
			pag.update(rsp)
			if err != nil {
				errs <- err
				return
			}
			for _, next := range batch {
				team := new(spec.Team)
				if err := team.FromGithub(next); err != nil {
					errs <- err
					return
				}
				model, err := spec.NewModel(spec.KindTeam)
				if err != nil {
					panic(err)
				}
				model.Metadata.ID = next.ID
				model.Metadata.Owner = org
				model.Metadata.Name = next.GetSlug()
				if next.Parent != nil {
					model.Metadata.Annotations["team/parent.id"] = strconv.FormatInt(next.Parent.GetID(), 10)
					model.Metadata.Annotations["team/parent.slug"] = next.Parent.GetSlug()
				}
				model.Spec = team
				res <- model
			}
		}
	}()
	return res, errs
}

func (g *Gitstrap) GetProtection(owner, repo, branch string) (*spec.Model, error) {
	if owner == "" {
		owner = g.me
	}
	ctx, cancel := g.newContext()
	defer cancel()
	res, _, err := g.gh.Repositories.GetBranchProtection(ctx, owner, repo, branch)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch protection rule: %w", err)
	}
	bp := new(spec.Protection)
	if err := bp.FromGithub(res); err != nil {
		return nil, fmt.Errorf("failed to parse protection spec: %w", err)
	}
	m, err := spec.NewModel(spec.KindProtection)
	if err != nil {
		panic(err)
	}
	m.Metadata.Name = branch
	m.Metadata.Repo = repo
	m.Metadata.Owner = owner
	m.Spec = bp
	return m, nil
}
