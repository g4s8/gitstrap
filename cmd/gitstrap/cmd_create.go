package main

import (
	"log"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/urfave/cli/v2"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create new resource",
	Action: cmdForEachModel(func(g *gitstrap.Gitstrap, m *spec.Model) error {
		if err := g.Create(m); err != nil {
			return err
		}
		log.Printf("Created: %s", m.Info())
		return nil
	}),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Resource specification file",
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force create, replace existing resource if exists",
		},
	},
}
