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
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "User name or organization owner of hooks repo",
				},
			},
		},
		{
			Name:   "teams",
			Usage:  "Get organization team",
			Action: cmdGetTeams,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "org",
					Usage:    "Organization name",
					Required: true,
				},
			},
		},
		{
			Name:   "protection",
			Usage:  "Get repository branch protection rules",
			Action: cmdGetProtections,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "owner",
					Usage:    "Repository name",
					Required: true,
				},
			},
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
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	stream, errs := g.GetHooks(c.String("owner"), name)
	enc := yaml.NewEncoder(os.Stdout)
	for {
		select {
		case h, ok := <-stream:
			if !ok {
				return nil
			}
			if err := enc.Encode(h); err != nil {
				return err
			}
		case err, ok := <-errs:
			if ok {
				return err
			}
		}
	}
}

func cmdGetTeams(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	stream, errs := g.GetTeams(ctx.String("org"))
	enc := yaml.NewEncoder(os.Stdout)
	for {
		select {
		case h, ok := <-stream:
			if !ok {
				return nil
			}
			if err := enc.Encode(h); err != nil {
				return err
			}
		case err, ok := <-errs:
			if ok {
				return err
			}
		}
	}
}

func cmdGetProtections(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	repo := ctx.Args().First()
	name := ctx.Args().Get(1)
	if name == "" || repo == "" {
		return fmt.Errorf("Requires repo and branch name argumentas")
	}
	s, err := g.GetProtection(ctx.String("owner"), repo, name)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(os.Stdout).Encode(s)
}
