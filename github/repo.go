package github

import (
	"errors"
	"fmt"
	"github.com/g4s8/gitstrap/context"
	"github.com/google/go-github/github"
	"log"
)

// Repo - GitHub repository
type Repo struct {
	Name          *string       `yaml:"name"`
	Description   *string       `yaml:"description"`
	Private       *bool         `yaml:"private,omitempty"`
	Hooks         RepoHooks     `yaml:"hooks"`
	Collaborators Collaborators `yaml:"collaborators"`
}

var errRepoWasSet = errors.New("Unexpected GitHub repository")
var errNoName = errors.New("Repo name was missed")

// Apply - apply configuration for new repo
func (r *Repo) Apply(ctx *context.Context) error {
	if ctx.GhRepo != nil {
		return errRepoWasSet
	}
	name, err := r.name()
	if err != nil {
		return err
	}
	ctx.GhRepo = &github.Repository{
		Name:        &name,
		Description: r.Description,
		Private:     r.Private,
	}
	o := ctx.Owner()
	log.Printf("applying configuration for repo %s/%s", o, name)
	var rep *github.Repository
	rep, resp, _ := ctx.Client.Repositories.Get(ctx.Sync, o, name)
	exists := resp.StatusCode == 200
	if exists {
		change := false
		if *rep.Description != *r.Description {
			*rep.Description = *r.Description
			change = true
		}
		if *rep.Private != r.IsPrivate() {
			*rep.Private = *r.Private
			change = true
		}
		if change {
			rep, _, err = ctx.Client.Repositories.Edit(ctx.Sync, o, name, rep)
			log.Printf("repository updated: %s", name)
			if err != nil {
				return fmt.Errorf("failed to update repo: %s", err)
			}
		} else {
			log.Printf("repository unchanged: %s", name)
		}
	} else {
		org := ""
		if optorg, hasOrg := ctx.Opt["org"]; hasOrg {
			org = optorg
		}
		rep = r.toGithub()
		rep, _, err = ctx.Client.Repositories.Create(ctx.Sync, org, rep)
		if err != nil {
			return fmt.Errorf("failed to create repo: %s", err)
		}
		log.Printf("repository created: %s", name)
	}
	ctx.GhRepo = rep
	if err := r.Hooks.Apply(ctx); err != nil {
		return fmt.Errorf("failed to apply webhooks: %s", err)
	}
	if err := r.Collaborators.Apply(ctx); err != nil {
		return fmt.Errorf("failed to apply collaborators: %s", err)
	}
	return nil
}

func (r *Repo) toGithub() *github.Repository {
	out := new(github.Repository)
	out.Name = r.Name
	out.Description = r.Description
	out.Private = r.Private
	return out
}

func (r *Repo) IsPrivate() bool {
	return r.Private != nil && *r.Private
}

func (r *Repo) name() (string, error) {
	if r.Name == nil || *r.Name == "" {
		return "", errNoName
	}
	return *r.Name, nil
}
