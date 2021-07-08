package main

import (
	"log"
	"os"
	"strconv"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/urfave/cli/v2"
)

var deleteCommand = &cli.Command{
	Name:    "delete",
	Aliases: []string{"remove", "del", "rm"},
	Usage:   "Delete resource",
	Action: cmdForEachModel(func(g *gitstrap.Gitstrap, m *spec.Model) error {
		if err := g.Delete(m); err != nil {
			return err
		}
		log.Printf("Deleted: %s", m.Info())
		return nil
	}),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Resource specification file",
		},
	},
	Subcommands: []*cli.Command{
		{
			Name:  "repo",
			Usage: "Delete repository",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
			},
			Action: cmdDeleteRepo,
		},
		{
			Name:  "readme",
			Usage: "Delete readme",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
			},
			Action: cmdDeleteReadme,
		},
		{
			Name:  "hook",
			Usage: "Delete webhook",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
				&cli.StringFlag{
					Name:  "repo",
					Usage: "Repository where hook is installed",
				},
			},
			Action: cmdDeleteHook,
		},
		{
			Name:  "team",
			Usage: "Delete team",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "org",
					Usage:    "Organization name",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "id",
					Usage: "Team ID",
				},
			},
			Action: cmdDeleteTeam,
		},
		{
			Name:  "protection",
			Usage: "Delete protection",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
			},
			Action: cmdDeleteProtection,
		},
	},
}

func cmdDeleteRepo(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	m, err := spec.NewModel(spec.KindRepo)
	if err != nil {
		return err
	}
	m.Metadata.Name = ctx.Args().First()
	m.Metadata.Owner = ctx.String("owner")
	if err := g.Delete(m); err != nil {
		return err
	}
	log.Printf("Deleted: %s", m.Info())
	return nil
}

func cmdDeleteReadme(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	m, err := spec.NewModel(spec.KindReadme)
	if err != nil {
		return err
	}
	spec := new(spec.Readme)
	spec.Selector.Repository = ctx.Args().First()
	m.Metadata.Owner = ctx.String("owner")
	m.Spec = spec
	if err := g.Delete(m); err != nil {
		return err
	}
	log.Printf("Deleted: %s", m.Info())
	return nil
}

func cmdDeleteHook(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	m, err := spec.NewModel(spec.KindHook)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(ctx.Args().First(), 10, 64)
	if err != nil {
		return err
	}
	m.Metadata.ID = &id
	spec := new(spec.Hook)
	if repo := ctx.String("repo"); repo != "" {
		spec.Selector.Repository = repo
		m.Metadata.Owner = ctx.String("owner")
	} else {
		spec.Selector.Organization = ctx.String("owner")
	}
	m.Spec = spec
	if err := g.Delete(m); err != nil {
		return err
	}
	log.Printf("Deleted: %s", m.Info())
	return nil
}

func cmdDeleteTeam(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	m, err := spec.NewModel(spec.KindTeam)
	if err != nil {
		return err
	}
	m.Metadata.Name = ctx.Args().First()
	m.Metadata.Owner = ctx.String("org")
	if id := ctx.String("id"); id != "" {
		iid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		m.Metadata.ID = &iid

	}
	if err := g.Delete(m); err != nil {
		return err
	}
	log.Printf("Deleted: %s", m.Info())
	return nil
}

func cmdDeleteProtection(ctx *cli.Context) error {
	token, err := resolveToken(ctx)
	if err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(ctx.Context, token, debug)
	if err != nil {
		return err
	}
	m, err := spec.NewModel(spec.KindProtection)
	if err != nil {
		return err
	}
	m.Metadata.Repo = ctx.Args().First()
	m.Metadata.Name = ctx.Args().Get(1)
	m.Metadata.Owner = ctx.String("owner")
	if err := g.Delete(m); err != nil {
		return err
	}
	log.Printf("Deleted: %s", m.Info())
	return nil
}
