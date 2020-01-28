package github

import (
	"errors"
	"github.com/g4s8/gitstrap/context"
	"log"
)

var errNoRepoName = errors.New("no repo name config")

// Delete GitHub repo
func (r *Repo) Delete(ctx *context.Context) error {
	name, err := r.name()
	if err != nil {
		return err
	}
	owner := ctx.Owner()
	if err != nil {
		return err
	}
	if _, err := ctx.Client.Repositories.Delete(ctx.Sync, owner, name); err != nil {
		return err
	}
	log.Printf("removed repository %s", name)
	return nil
}
