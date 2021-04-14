package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func TeamExist(cli *gh.Client, ctx context.Context, org, slug string) (bool, error) {
	_, _, err := cli.Teams.GetTeamBySlug(ctx, org, slug)
	return resolveResponseByErr(err)
}
