package main

import (
	"errors"
	"fmt"
	"github.com/g4s8/gitstrap"
	"github.com/g4s8/gitstrap/config"
	"github.com/g4s8/gitstrap/context"
	"github.com/g4s8/gopwd"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	buildVersion string
	buildCommit  string
	buildDate    string
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	app := cli.App{
		Name:        "gitstrap",
		Description: "CLI tool to manage GitHub repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "GitHub API token with repo access",
			},
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get repository configuration",
				Action:  cmdGet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "names",
						Usage:   "output format: names|yaml",
					},
					&cli.StringFlag{
						Name:  "org",
						Usage: "GitHub organization owner instead of user",
					},
				},
			},
			&cli.Command{
				Name:   "apply",
				Usage:  "apply gitstrap configuration",
				Action: cmdApply,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "gitstrap yaml config to apply",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "org",
						Usage: "GitHub organization owner instead of user",
					},
				},
			},
			&cli.Command{
				Name:   "delete",
				Usage:  "delete gitstrap repository",
				Action: cmdDelete,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "gitstrap yaml config",
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "accept",
						Usage:    "accept delete",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "org",
						Usage: "GitHub organization owner instead of user",
					},
				},
			},
			&cli.Command{
				Name:        "upgrade",
				Description: "upgrade config to latest version",
				Action:      cmdUpgrade,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "gitstrap yaml config",
						Required: true,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func cmdUpgrade(c *cli.Context) error {
	f := c.String("file")
	cfg := new(config.Config)
	if err := cfg.ParseFile(f); err != nil {
		return err
	}
	y, err := cfg.ToYaml()
	if err != nil {
		return err
	}
	fmt.Println(y)
	return nil
}

func cmdDelete(c *cli.Context) error {
	token := c.String("token")
	var err error
	if token, err = getToken(token); err != nil {
		return err
	}
	f := c.String("file")
	cfg := new(config.Config)
	if err := cfg.ParseFile(f); err != nil {
		return err
	}
	pwd, err := gopwd.Abs()
	if err != nil {
		return err
	}
	ctx := context.New(token, pwd, false)
	org := c.String("org")
	if org != "" {
		ctx.Opt["org"] = org
	}
	if err := gitstrap.Delete(ctx, cfg); err != nil {
		return err
	}
	return nil
}

func cmdApply(c *cli.Context) error {
	token := c.String("token")
	var err error
	if token, err = getToken(token); err != nil {
		return err
	}
	f := c.String("file")
	cfg := new(config.Config)
	if err := cfg.ParseFile(f); err != nil {
		return err
	}
	pwd, err := gopwd.Abs()
	if err != nil {
		return err
	}
	ctx := context.New(token, pwd, false)
	org := c.String("org")
	if org != "" {
		ctx.Opt["org"] = org
	}
	if err := gitstrap.Apply(ctx, cfg); err != nil {
		return err
	}
	return nil
}

func cmdGet(c *cli.Context) error {
	token := c.String("token")
	var err error
	if token, err = getToken(token); err != nil {
		return err
	}
	name := c.Args().Get(0)
	if name == "" {
		return errors.New("name argument required")
	}
	pwd, err := gopwd.Abs()
	if err != nil {
		return err
	}
	ctx := context.New(token, pwd, false)
	org := c.String("org")
	if org != "" {
		ctx.Opt["org"] = org
	}
	cfg := new(config.Config)
	if err := gitstrap.Get(ctx, cfg, name); err != nil {
		return err
	}
	fout := c.String("output")
	if fout == "yaml" {
		yaml, err := cfg.ToYaml()
		if err != nil {
			return err
		}
		fmt.Print(yaml)
	} else if fout == "names" {
		fmt.Printf("name\tdescription\n%s\t%s\n", *cfg.Github.Name, *cfg.Github.Description)
	} else {
		return fmt.Errorf("unsupported output format: %s", fout)
	}
	return nil
}

func getToken(token string) (string, error) {
	if token != "" {
		return token, nil
	}
	file := os.Getenv("HOME") + "/.config/gitstrap/github_token.txt"
	if token, err := ioutil.ReadFile(file); err == nil {
		return string(token), nil
	}
	return "", fmt.Errorf("GitHub token neither given as a flag, nor found in %s", file)
}

func fatal(err error) {
	if _, xerr := fmt.Fprintf(os.Stderr, "%s\n", err); xerr != nil {
		fmt.Printf("Failed to print error: %s", xerr)
	}
	os.Exit(1)
}
