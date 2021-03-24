package gitstrap

import (
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/google/go-github/v33/github"
)

// Get resource
func (g *Gitstrap) Get(name string, owner string, format spec.ModelFormat) (<-chan view.Printable, <-chan error) {
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
		cls := make([]*github.User, 0)
		copts := &github.ListCollaboratorsOptions{Affiliation: "direct"}
	PAGINATION:
		batch, rsp, err := g.gh.Repositories.ListCollaborators(ctx, owner, name, copts)
		cls = append(cls, batch...)
		if rsp.NextPage < rsp.LastPage {
			copts.Page = rsp.NextPage
			goto PAGINATION
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
		repo.Collaborators.FromUsers(cls)
		model.Spec = repo
		if p, err := format.ToView(model); err != nil {
			errs <- err
		} else {
			res <- p
		}
	}()
	return res, errs
}
