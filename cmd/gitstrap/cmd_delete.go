package main

import (
	"log"

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
}
