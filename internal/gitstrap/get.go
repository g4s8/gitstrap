package gitstrap

import (
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/google/go-github/v33/github"
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

const (
	hooksPageSize = 10
)

func (g *Gitstrap) GetRepoHooks(owner, name string) (<-chan *spec.Model, <-chan error) {
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
		for {
			ghooks, rsp, err := g.gh.Repositories.ListHooks(ctx, owner, name, opts)
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
				hook.Selector.Repository = name
				if err := hook.FromGithub(gh); err != nil {
					errs <- err
					return
				}
				s.Spec = hook
				out <- s
			}
			if opts.Page == rsp.LastPage {
				break
			}
			opts.Page = rsp.NextPage
		}

	}()
	return out, errs
}
