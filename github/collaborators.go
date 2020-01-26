package github

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
	"log"
)

// Collaborators - GitHub repository collaborators
type Collaborators []string

func (cls Collaborators) Apply(ctx *context.Context) error {
	own := repoOwner(ctx.GhRepo)
	name := *ctx.GhRepo.Name
	users, _, err := ctx.Client.Repositories.ListCollaborators(ctx.Sync, own, name, nil)
	if err != nil {
		return fmt.Errorf("failed to list collabolators: %s", err)
	}
	add := make([]string, 0)
	del := make([]string, 0)
	for _, u := range users {
		if !cls.contains(*u.Login) {
			del = append(del, *u.Login)
		}
	}
	for _, c := range cls {
		has := false
		for _, u := range users {
			if u.GetLogin() == c {
				has = true
				break
			}
		}
		if !has {
			add = append(add, c)
		}
	}
	for _, u := range del {
		if u == own {
			continue
		}
		log.Printf("removing collaborator %s", u)
		if _, err := ctx.Client.Repositories.RemoveCollaborator(ctx.Sync, own, name, u); err != nil {
			return fmt.Errorf("failed to delete collaborator: %s", err)
		}
	}
	for _, u := range add {
		log.Printf("adding collaborator %s", u)
		if _, err := ctx.Client.Repositories.AddCollaborator(ctx.Sync, own, name, u, nil); err != nil {
			return fmt.Errorf("failed to add collaborator: %s", err)
		}
	}
	return nil
}

func (cls Collaborators) contains(name string) bool {
	for _, c := range cls {
		if c == name {
			return true
		}
	}
	return false
}
