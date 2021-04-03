package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func OrgExist(cli *gh.Client, ctx context.Context, name string) (bool, error) {
	_, rsp, err := cli.Organizations.Get(ctx, name)
	if rsp.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
