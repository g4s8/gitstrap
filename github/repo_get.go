package github

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
)

// Get GitHub repository config.
func (r *Repo) Get(ctx *context.Context, name string) error {
	owner := ctx.Owner()
	rep, rsp, err := ctx.Client.Repositories.Get(ctx.Sync, owner, name)
	if rsp.StatusCode == 404 {
		return fmt.Errorf("repository %s/%s not found", owner, name)
	}
	if err != nil {
		return fmt.Errorf("failed to get repository: %s", err)
	}
	r.Name = rep.Name
	r.Description = rep.Description
	r.Private = rep.Private

	users, _, err := ctx.Client.Repositories.ListCollaborators(ctx.Sync,
		owner, name, nil)
	if err != nil {
		return fmt.Errorf("failed to get collaborators: %s", err)
	}
	cls := make([]string, len(users), len(users))
	for i, u := range users {
		cls[i] = *u.Login
	}
	r.Collaborators = Collaborators(cls)
	ghks, _, err := ctx.Client.Repositories.ListHooks(ctx.Sync, owner, name, nil)
	if err != nil {
		return err
	}
	hooks := make([]RepoHook, len(ghks), len(ghks))
	for i, gh := range ghks {
		hooks[i] = RepoHook{
			URL:    gh.Config["url"].(string),
			Type:   gh.Config["content_type"].(string),
			Active: gh.Active,
			Events: gh.Events,
		}
	}
	r.Hooks = hooks

	return nil
}
