package main

import (
	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
)

var applyCommand = &cli.Command{
	Name:  "apply",
	Usage: "Apply new specficiation",
	Action: cmdForEachModel(func(g *gitstrap.Gitstrap, m *spec.Model) error {
		rs, errs := g.Apply(m)
		return view.RenderOn(view.Console, rs, errs)
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
