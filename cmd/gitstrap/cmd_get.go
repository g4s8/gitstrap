package main

import (
	"fmt"
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var getCommand = &cli.Command{
	Name:    "get",
	Aliases: []string{"g"},
	Usage:   "Get resource",
	Subcommands: []*cli.Command{
		{
			Name:   "repo",
			Usage:  "Get repository",
			Action: cmdGetRepo,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Get repositories of another user or organization",
				},
			},
		},
		{
			Name:   "org",
			Usage:  "Get organization",
			Action: cmdGetOrg,
		},
		{
			Name:   "hooks",
			Usage:  "Get webhooks configurations",
			Action: cmdGetHooks,
		},
	},
}

func cmdGetRepo(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}
	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("Requires repository name argument")
	}
	owner := c.String("owner")
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	repo, err := g.GetRepo(name, owner)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(os.Stdout).Encode(repo)
}

func cmdGetOrg(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}
	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("Requires repository name argument")
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	org, err := g.GetOrg(name)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(os.Stdout).Encode(org)
}

func cmdGetHooks(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}
	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("Requires repository name argument")
	}
	debug := os.Getenv("DEBUG") != ""
	_, err = gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	return nil
}
