package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func OrgExist(cli *gh.Client, ctx context.Context, name string) (bool, error) {
	_, _, err := cli.Organizations.Get(ctx, name)
	return resolveResponseByErr(err)
}
