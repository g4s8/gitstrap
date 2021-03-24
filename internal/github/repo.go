package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func RepoExist(cli *gh.Client, ctx context.Context, owner, name string) (bool, error) {
	_, rsp, err := cli.Repositories.Get(ctx, owner, name)
	if rsp.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
