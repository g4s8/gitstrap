package context

import "errors"

// ResolveOwner resolves GitHub repo owner.
// It's either current GitHub user login resolved from
// API token or organization login provided as `-org`
// option.
func (ctx *Context) ResolveOwner() error {
	if org, hasOrg := ctx.Opt["org"]; hasOrg {
		ctx.setOwner(org)
		return nil
	}
	me, _, err := ctx.Client.Users.Get(ctx.Sync, "")
	if err != nil {
		return err
	}
	ctx.setOwner(*me.Login)
	return nil
}

var errOwnerAlreadySet = errors.New("owner was set already")

func (ctx *Context) setOwner(val string) {
	if ctx.owner != nil {
		panic(errOwnerAlreadySet)
	}
	ctx.owner = new(string)
	*ctx.owner = val
}

var errOwnerNotResolved = errors.New("owner was not resolved")

// Owner of GitHub repository
func (ctx *Context) Owner() string {
	if ctx.owner == nil {
		panic(errOwnerNotResolved)
	}
	return *ctx.owner
}

var errRepoWasNotSet = errors.New("github repo was not set in the context")

func (ctx *Context) Name() string {
	if ctx.GhRepo == nil || ctx.GhRepo.Name == nil {
		panic(errRepoWasNotSet)
	}
	return *ctx.GhRepo.Name
}
