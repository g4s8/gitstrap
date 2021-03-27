package main

import (
	"fmt"
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
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
	format := spec.MfYaml
	owner := c.String("owner")
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	repo, errs := g.GetRepo(name, owner, format)
	if err := view.RenderOn(view.Console, repo, errs); err != nil {
		fatal(err)
	}
	return nil
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
	format := spec.MfYaml
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	repo, errs := g.GetOrg(name, format)
	if err := view.RenderOn(view.Console, repo, errs); err != nil {
		fatal(err)
	}
	return nil
}
