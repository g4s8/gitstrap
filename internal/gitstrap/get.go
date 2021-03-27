package gitstrap

import (
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
)

// GetRepo repository resource
func (g *Gitstrap) GetRepo(name string, owner string, format spec.ModelFormat) (<-chan view.Printable, <-chan error) {
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	if owner == "" {
		owner = g.me
	}
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		r, _, err := g.gh.Repositories.Get(ctx, owner, name)
		if err != nil {
			errs <- err
			return
		}
		model := new(spec.Model)
		model.Kind = spec.KindRepo
		model.Version = spec.Version
		model.Metadata = new(spec.Metadata)
		model.Metadata.Name = name
		model.Metadata.ID = r.ID
		if owner != "" {
			model.Metadata.Owner = owner
		} else {
			model.Metadata.Owner = g.me
		}
		repo := new(spec.Repo)
		repo.FromGithub(r)
		model.Spec = repo
		if p, err := format.ToView(model); err != nil {
			errs <- err
		} else {
			res <- p
		}
	}()
	return res, errs
}

func (g *Gitstrap) GetOrg(name string, format spec.ModelFormat) (<-chan view.Printable, <-chan error) {
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		r, _, err := g.gh.Organizations.Get(ctx, name)
		if err != nil {
			errs <- err
			return
		}
		model := new(spec.Model)
		model.Kind = spec.KindOrg
		model.Version = spec.Version
		model.Metadata = new(spec.Metadata)
		model.Metadata.Name = name
		model.Metadata.ID = r.ID
		org := new(spec.Org)
		org.FromGithub(r)
		model.Spec = org
		v, err := format.ToView(model)
		if err != nil {
			errs <- err
		} else {
			res <- v
		}
	}()
	return res, errs
}
