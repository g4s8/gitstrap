package gitstrap

import (
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/google/go-github/v33/github"
)

type errUnsupportModelKind struct {
	kind spec.Kind
}

func (e *errUnsupportModelKind) Error() string {
	return fmt.Sprintf("Unsupported model kind: `%s`", e.kind)
}

func (g *Gitstrap) Create(m *spec.Model) error {
	switch m.Kind {
	case spec.KindRepo:
		return g.createRepo(m)
	case spec.KindReadme:
		return g.createReadme(m)
	case spec.KindHook:
		return g.createHook(m)
	default:
		return &errUnsupportModelKind{m.Kind}
	}
}

func (g *Gitstrap) createRepo(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	meta := m.Metadata
	repo := new(spec.Repo)
	if err := m.GetSpec(repo); err != nil {
		return err
	}
	if repo.DefaultBranch == "" {
		repo.DefaultBranch = "master"
	}
	owner := g.getOwner(m)
	fn := fmt.Sprintf("%s/%s", owner, meta.Name)
	grepo := new(github.Repository)
	if err := repo.ToGithub(grepo); err != nil {
		return err
	}
	grepo.Name = &meta.Name
	grepo.FullName = &fn
	org := g.resolveOrg(m)
	r, _, err := g.gh.Repositories.Create(ctx, org, grepo)
	if err != nil {
		return err
	}
	m.Metadata.FromGithubRepo(r)
	repo.FromGithub(r)
	m.Spec = repo
	return nil
}

type errReadmeExists struct {
	owner, repo string
}

func (e *errReadmeExists) Error() string {
	return fmt.Sprintf("README.md already exists in %s/%s (try --force for replacing it)", e.owner, e.repo)
}

type errReadmeNotFile struct {
	rtype string
}

func (e *errReadmeNotFile) Error() string {
	return fmt.Sprintf("README is no a file: `%s`", e.rtype)
}

func (g *Gitstrap) createReadme(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	meta := m.Metadata
	spec := new(spec.Readme)
	if err := m.GetSpec(spec); err != nil {
		return err
	}
	owner := m.Metadata.Owner
	if owner == "" {
		owner = g.me
	}
	repo, _, err := g.gh.Repositories.Get(ctx, owner, spec.Selector.Repository)
	if err != nil {
		return err
	}
	msg := "Updated README.md"
	if cm, ok := meta.Annotations["commitMessage"]; ok {
		msg = cm
	}
	opts := &github.RepositoryContentFileOptions{
		Content: []byte(spec.String()),
		Message: &msg,
	}
	if meta.Annotations["force"] == "true" {
		getopts := &github.RepositoryContentGetOptions{}
		cnt, _, rsp, err := g.gh.Repositories.GetContents(ctx, owner, repo.GetName(), "README.md", getopts)
		if rsp.StatusCode == 404 {
			goto SKIP_GET
		}
		if err != nil {
			return err
		}
		if *cnt.Type != "file" {
			return &errReadmeNotFile{*cnt.Type}
		}
		opts.SHA = cnt.SHA
	SKIP_GET:
	}
	_, rsp, err := g.gh.Repositories.UpdateFile(ctx, owner, repo.GetName(), "README.md", opts)
	if err != nil {
		if rsp.StatusCode == 422 && opts.SHA == nil {
			return &errReadmeExists{owner, repo.GetName()}
		}
		return err
	}
	return nil
}

func (g *Gitstrap) createHook(m *spec.Model) error {
	ctx, cancel := g.newContext()
	defer cancel()
	owner := g.getOwner(m)
	hook := new(spec.Hook)
	if err := m.GetSpec(hook); err != nil {
		return err
	}
	ghook := new(github.Hook)
	if err := hook.ToGithub(ghook); err != nil {
		return err
	}
	var err error
	if hook.Selector.Repository != "" {
		ghook, _, err = g.gh.Repositories.CreateHook(ctx, owner, hook.Selector.Repository, ghook)
	} else if hook.Selector.Organization != "" {
		ghook, _, err = g.gh.Organizations.CreateHook(ctx, hook.Selector.Organization, ghook)
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
