package main

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var initCommand = &cli.Command{
	Name:  "init",
	Usage: "Generate stub specification file",
	Subcommands: []*cli.Command{
		{
			Name:  "repo",
			Usage: "Generate repo stub",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
			},
			Action: initCmd(func(ctx *cli.Context) (*spec.Model, error) {
				m, err := spec.NewModel(spec.KindRepo)
				if err != nil {
					return nil, err
				}
				if n := ctx.Args().First(); n != "" {
					m.Metadata.Name = n
				} else {
					m.Metadata.Name = "repo"
				}
				m.Metadata.Owner = ctx.String("owner")
				spec := *new(spec.Repo)
				m.Spec = spec
				return m, nil
			}),
		},
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "full",
			Value:   false,
			Usage:   "Init full spec with empty and default fields",
			Aliases: []string{"f"},
		},
	},
}

func initCmd(model func(*cli.Context) (*spec.Model, error)) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		m, err := model(ctx)
		if err != nil {
			return err
		}
		s := new(spec.Repo)
		if err := defaults.Set(s); err != nil {
			panic(err)
		}
		if ctx.Bool("full") {
			spe, _ := spec.RemoveTagsOmitempty(*s)
			m.Spec = spe
		} else {
			m.Spec = s
		}
		return yaml.NewEncoder(os.Stdout).Encode(m)
	}
}
