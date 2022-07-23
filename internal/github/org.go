package github

import (
	"context"

	gh "github.com/google/go-github/v38/github"
)

func OrgExist(cli *gh.Client, ctx context.Context, name string) (bool, error) {
	_, _, err := cli.Organizations.Get(ctx, name)
	if isNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func GetOrgIdByName(cli *gh.Client, ctx context.Context, name string) (int64, error) {
	org, _, err := cli.Organizations.Get(ctx, name)
	if err != nil {
		return 0, err
	}
	return *org.ID, nil
}
