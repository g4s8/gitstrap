package main

import (
	"fmt"
	"os"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/urfave/cli/v2"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls", "lst"},
	Usage:   "List resources",
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
			Action: cmdListRepo,
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
	lst, errs := g.ListRepos(filter, owner)
	out := os.Stdout
	var done bool
	for !done {
		select {
		case next, ok := <-lst:
			if !ok {
				done = true
				break
			}
			if _, err := next.WriteTo(out); err != nil {
				return err
			}
			fmt.Print(out, "+\n")
		case err, ok := <-errs:
			if !ok {
				done = true
				break
			}
			return err
		}
	}
	return nil
}
