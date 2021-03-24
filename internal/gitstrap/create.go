package gitstrap

import (
	"context"
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/google/go-github/v33/github"
)

type errUnsupportModelKind struct {
	kind string
}

func (e *errUnsupportModelKind) Error() string {
	return fmt.Sprintf("Unsupported model kind: `%s`", e.kind)
}

func (g *Gitstrap) Create(m *spec.Model) (<-chan view.Printable, <-chan error) {
	res := make(chan view.Printable)
	errs := make(chan error)
	ctx, cancel := g.newContext()
	go func() {
		defer close(res)
		defer close(errs)
		defer cancel()
		switch m.Kind {
		case spec.KindRepo:
			spec := m.Spec.(*spec.Repo)
			r, err := g.createRepo(ctx, spec, m.Metadata)
			if err != nil {
				errs <- err
			} else {
				res <- r
			}
		case spec.KindReadme:
			spec := m.Spec.(*spec.Readme)
			r, err := g.createReadme(ctx, spec, m.Metadata)
			if err != nil {
				errs <- err
			} else {
				res <- r
			}
		default:
			errs <- &errUnsupportModelKind{m.Kind}
		}
	}()
	return res, errs
}

type createRepoResult struct {
	repo *github.Repository
}

func (cr *createRepoResult) PrintOn(p view.Printer) {
	p.Print(fmt.Sprintf("Repository [%d] %s created", cr.repo.GetID(), cr.repo.GetFullName()))
}

func (g *Gitstrap) createRepo(ctx context.Context, repo *spec.Repo, meta *spec.Metadata) (view.Printable, error) {
	grepo := new(github.Repository)
	if err := repo.ToGithub(grepo); err != nil {
		return nil, err
	}
	grepo.Name = &meta.Name
	fn := fmt.Sprintf("%s/%s", meta.Owner, meta.Name)
	grepo.FullName = &fn
	owner := meta.Owner
	if owner == "" || owner == g.me {
		owner = ""
	}
	r, _, err := g.gh.Repositories.Create(ctx, owner, grepo)
	if err != nil {
		return nil, err
	}
	return &createRepoResult{r}, nil
}

type createReadmeResult struct {
	*github.RepositoryContentResponse
}

func (r *createReadmeResult) PrintOn(p view.Printer) {
	p.Print(fmt.Sprintf("README created with %s", r.GetSHA()))
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

func (g *Gitstrap) createReadme(ctx context.Context, spec *spec.Readme, meta *spec.Metadata) (view.Printable, error) {
	owner := spec.Selector.Owner
	if owner == "" {
		owner = g.me
	}
	repo, _, err := g.gh.Repositories.Get(ctx, owner, spec.Selector.Repository)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if *cnt.Type != "file" {
			return nil, &errReadmeNotFile{*cnt.Type}
		}
		opts.SHA = cnt.SHA
	SKIP_GET:
	}
	rs, rsp, err := g.gh.Repositories.UpdateFile(ctx, owner, repo.GetName(), "README.md", opts)
	if err != nil {
		if rsp.StatusCode == 422 && opts.SHA == nil {
			return nil, &errReadmeExists{owner, repo.GetName()}
		}
		return nil, err
	}
	return &createReadmeResult{rs}, nil
}
