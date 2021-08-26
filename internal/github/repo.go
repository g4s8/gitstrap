package github

import (
	"context"

	gh "github.com/google/go-github/v38/github"
)

func RepoExist(cli *gh.Client, ctx context.Context, owner, name string) (bool, error) {
	_, _, err := cli.Repositories.Get(ctx, owner, name)
	return resolveResponseByErr(err)
}
