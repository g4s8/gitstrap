package main

import (
	"log"
	"os"

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
			Name: "repo",
			Usage: "Delete repository",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "owner",
					Usage: "Repository owner or organization name",
				},
			},
			Action: cmdDeleteRepo,
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
