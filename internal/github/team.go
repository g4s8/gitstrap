package github

import (
	"context"

	gh "github.com/google/go-github/v38/github"
)

func TeamExistBySlug(cli *gh.Client, ctx context.Context, org, slug string) (bool, error) {
	_, _, err := cli.Teams.GetTeamBySlug(ctx, org, slug)
	if isNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func TeamExistByID(cli *gh.Client, ctx context.Context, org, ID int64) (bool, error) {
	_, _, err := cli.Teams.GetTeamByID(ctx, org, ID)
	if isNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
