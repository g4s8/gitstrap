package main

import (
	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
)

var deleteCommand = &cli.Command{
	Name:    "delete",
	Aliases: []string{"remove", "del", "rm"},
	Usage:   "Delete resource",
	Action: cmdForEachModel(func(g *gitstrap.Gitstrap, m *spec.Model) error {
		rs, errs := g.Delete(m)
		return view.RenderOn(view.Console, rs, errs)
	}),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Resource specification file",
		},
	},
}
