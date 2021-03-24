package main

import (
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create new resource",
	Action:  cmdCreate,
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

func cmdCreate(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}

	model := new(spec.Model)
	if err := model.FromFile(c.String("file")); err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	if c.Bool("force") {
		model.Metadata.Annotations["force"] = "true"
	}
	rs, errs := g.Create(model)
	if err := view.RenderOn(view.Console, rs, errs); err != nil {
		fatal(err)
	}
	return nil
}
