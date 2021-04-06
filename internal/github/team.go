package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func TeamExist(cli *gh.Client, ctx context.Context, org, slug string) (bool, error) {
	_, rsp, err := cli.Teams.GetTeamBySlug(ctx, org, slug)
	if err != nil {
		return false, err
	}
	if rsp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}
