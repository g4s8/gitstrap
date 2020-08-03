package github

import (
	"github.com/g4s8/gitstrap/context"
	"log"
)

// Delete GitHub repo
func (r *Repo) Delete(ctx *context.Context) error {
	name, err := r.name()
	if err != nil {
		return err
	}
	owner := ctx.Owner()
	if _, err := ctx.Client.Repositories.Delete(ctx.Sync, owner, name); err != nil {
		return err
	}
	log.Printf("removed repository %s", name)
	return nil
}
