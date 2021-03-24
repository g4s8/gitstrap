package main

import (
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/view"
	"github.com/urfave/cli/v2"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls", "lst"},
	Usage:   "List resources",
	Action:  cmdListRepo,
	Subcommands: []*cli.Command{
		{
			Name:  "repo",
			Usage: "List repositories",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "List repositories of another user or organization",
				},
				&cli.BoolFlag{
					Name:  "forks",
					Usage: "Filter only fork repositories",
				},
				&cli.BoolFlag{
					Name:  "no-forks",
					Usage: "Filter out fork repositories",
				},
				&cli.IntFlag{
					Name:  "stars-gt",
					Usage: "Filter by stars greater than value",
				},
				&cli.IntFlag{
					Name:  "stars-lt",
					Usage: "Filter by stars less than value",
				},
			},
		},
	},
}

func cmdListRepo(c *cli.Context) error {
	token, err := resolveToken(c)
	if err != nil {
		return err
	}
	owner := c.String("owner")
	debug := os.Getenv("DEBUG") != ""
	g, err := gitstrap.New(c.Context, token, debug)
	if err != nil {
		return err
	}
	filter := gitstrap.LfNop
	if c.Bool("forks") {
		filter = gitstrap.LfForks(filter, true)
	}
	if c.Bool("no-forks") {
		filter = gitstrap.LfForks(filter, false)
	}
	if gt := c.Int("stars-gt"); gt > 0 {
		filter = gitstrap.LfStars(filter, gitstrap.LfStarsGt(gt))
	}
	if lt := c.Int("stars-lt"); lt > 0 {
		filter = gitstrap.LfStars(filter, gitstrap.LfStarsLt(lt))
	}
	lst, errs := g.List(filter, owner)
	if err := view.RenderOn(view.Console, lst, errs); err != nil {
		fatal(err)
	}
	return nil
}
