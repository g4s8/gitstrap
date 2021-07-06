package main

import (
	"fmt"
	"os"
	"reflect"

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
	errs := make(chan error)
	defer close(errs)
	lst := g.ListRepos(filter, owner, errs)
	cases := make([]reflect.SelectCase, 2)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(lst)}
	cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(errs)}
	out := os.Stdout
	for {
		choosen, val, ok := reflect.Select(cases)
		if !ok {
			break
		}
		if choosen == 0 {
			next := val.Interface().(*gitstrap.RepoInfo)
			if _, err := next.WriteTo(out); err != nil {
				return err
			}
			fmt.Fprint(out, "\n")
		} else {
			return val.Interface().(error)
		}
	}
	return nil
}
