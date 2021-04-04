package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func OrgExist(cli *gh.Client, ctx context.Context, name string) (bool, error) {
	_, rsp, err := cli.Organizations.Get(ctx, name)
	if err != nil {
		return false, err
	}
	if rsp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}
