package github

import (
	"context"

	gh "github.com/google/go-github/v33/github"
)

func TeamExistBySlug(cli *gh.Client, ctx context.Context, org, slug string) (bool, error) {
	_, _, err := cli.Teams.GetTeamBySlug(ctx, org, slug)
	return resolveResponseByErr(err)
}

func TeamExistByID(cli *gh.Client, ctx context.Context, org, ID int64) (bool, error) {
	_, _, err := cli.Teams.GetTeamByID(ctx, org, ID)
	return resolveResponseByErr(err)
}
