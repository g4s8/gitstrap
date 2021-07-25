package main

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	stubOrg   = "exampleOrg"
	stubRepo  = "exampleRepo"
	stubOwner = "exampleOwner"
	stubTeam  = "exampleTeam"
)

var initCommand = &cli.Command{
	Name:  "init",
	Usage: "Generate stub specification file",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "full",
			Value:   false,
			Usage:   "Init full spec with empty and default fields",
			Aliases: []string{"f"},
		},
	},
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
					m.Metadata.Name = stubRepo
				}
				m.Metadata.Owner = ctx.String("owner")
				spec := new(spec.Repo)
				m.Spec = spec
				return m, nil
			}),
		},
		{
			Name:  "org",
			Usage: "Generate org stub",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "id",
					Usage: "Organization ID",
				},
			},
			Action: initCmd(func(ctx *cli.Context) (*spec.Model, error) {
				m, err := spec.NewModel(spec.KindOrg)
				if err != nil {
					return nil, err
				}
				if n := ctx.Args().First(); n != "" {
					m.Metadata.Name = n
				} else {
					m.Metadata.Name = stubOrg
				}
				if ID := ctx.Int64("id"); ID != 0 {
					m.Metadata.ID = &ID
				}
				spec := new(spec.Org)
				spec.Name = m.Metadata.Name
				m.Spec = spec
				return m, nil
			}),
		},
		{
			Name:  "hook",
			Usage: "Generate hook stub",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "id",
					Usage: "Hook ID, required on update",
				},
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Repository owner or organization name",
				},
				&cli.StringFlag{
					Name:  "repo",
					Usage: "Name of repository for this hook",
				},
			},
			Action: initCmd(func(ctx *cli.Context) (*spec.Model, error) {
				m, err := spec.NewModel(spec.KindHook)
				if ID := ctx.Int64("id"); ID != 0 {
					m.Metadata.ID = &ID
				}
				spec := new(spec.Hook)
				if err != nil {
					return nil, err
				}
				owner := ctx.String("owner")
				if repo := ctx.String("repo"); repo != "" {
					spec.Selector.Repository = repo
					m.Metadata.Owner = owner
				} else if owner != "" {
					spec.Selector.Organization = owner
				} else {
					spec.Selector.Organization = stubOrg
				}
				m.Spec = spec
				return m, nil
			}),
		},
		{
			Name:  "readme",
			Usage: "Generate readme stub",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repo",
					Usage: "Name of repository where this readme will be created",
				},
			},
			Action: initCmd(func(ctx *cli.Context) (*spec.Model, error) {
				m, err := spec.NewModel(spec.KindReadme)
				spec := new(spec.Readme)
				if err != nil {
					return nil, err
				}
				if repo := ctx.String("repo"); repo != "" {
					spec.Selector.Repository = repo
				} else {
					spec.Selector.Repository = stubRepo
				}
				m.Spec = spec
				return m, nil
			}),
		},
		{
			Name:  "team",
			Usage: "Generate team stub",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "slug",
					Usage: "Team slug",
				},
				&cli.StringFlag{
					Name:  "owner",
					Usage: "Organization to which the team belongs",
				},
				&cli.Int64Flag{
					Name:  "id",
					Usage: "Team ID",
				},
			},
			Action: initCmd(func(ctx *cli.Context) (*spec.Model, error) {
				m, err := spec.NewModel(spec.KindTeam)
				spec := new(spec.Team)
				if err != nil {
					return nil, err
				}
				if ID := ctx.Int64("id"); ID != 0 {
					m.Metadata.ID = &ID
				}
				if owner := ctx.String("owner"); owner != "" {
					m.Metadata.Owner = owner
				} else {
					m.Metadata.Owner = stubOrg
				}
				if slug := ctx.String("slug"); slug != "" {
					m.Metadata.Name = slug
					spec.Name = slug
				}
				m.Spec = spec
				return m, nil
			}),
		},
	},
}

func initCmd(model func(*cli.Context) (*spec.Model, error)) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		m, err := model(ctx)
		if err != nil {
			return err
		}
		if err := defaults.Set(m.Spec); err != nil {
			return err
		}
		if ctx.Bool("full") {
			s, err := spec.RemoveTagsOmitempty(m.Spec)
			if err != nil {
				return err
			}
			m.Spec = s
		}
		return yaml.NewEncoder(os.Stdout).Encode(m)
	}
}