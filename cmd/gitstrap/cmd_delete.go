package main

import (
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

var deleteCommand = &cli.Command{
	Name:    "delete",
	Aliases: []string{"remove", "del", "rm"},
	Usage:   "Delete resource",
	Action:  cmdDelete,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Resource specification file",
		},
	},
}

func cmdDelete(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}

	fn, _ := filepath.Abs(c.String("file"))
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	model := new(spec.Model)
	if err := yaml.Unmarshal(data, model); err != nil {
		return err
	}
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	rs, errs := g.Delete(model)
	if err := view.RenderOn(view.Console, rs, errs); err != nil {
		fatal(err)
	}
	return nil
}
