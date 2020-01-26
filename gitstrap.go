package gitstrap

import (
	"bufio"
	"fmt"
	"github.com/g4s8/gitstrap/config"
	"github.com/g4s8/gitstrap/context"
	"github.com/google/go-github/github"
	"os"
	"os/exec"
	"strings"
)

// Apply gitstrap configuration
func Apply(ctx *context.Context, cfg *config.Config) error {
	if err := ctx.ResolveOwner(); err != nil {
		return err
	}
	if err := cfg.Github.Apply(ctx); err != nil {
		return err
	}
	if err := cfg.Templates.Apply(ctx); err != nil {
		return err
	}
	if err := cfg.Ext.Apply(ctx); err != nil {
		return err
	}
	return nil
}

// Get configuration from existing repo
func Get(ctx *context.Context, cfg *config.Config, name string) error {
	if err := ctx.ResolveOwner(); err != nil {
		return err
	}
	if err := cfg.Get(ctx, name); err != nil {
		return err
	}
	return nil
}

// Delete gitstrap
func Delete(ctx *context.Context, cfg *config.Config) error {
	if err := ctx.ResolveOwner(); err != nil {
		return err
	}
	if err := cfg.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func gitSync(repo *github.Repository) error {
	if err := exec.Command("git", "init", ".").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "remote", "add", "origin", *repo.SSHURL).Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "fetch", "origin").Run(); err != nil {
		return err
	}
	fmt.Printf("Github repository %s has been fetched\n", *repo.SSHURL)
	return nil
}

func gitPush(repo *github.Repository) error {
	if prompt("Templates has been applied. Do you want to commit & push?") {
		if err := exec.Command("git", "add", ".").Run(); err != nil {
			return err
		}
		if err := exec.Command("git", "commit",
			"-m", "[gitstrap] bootstrap repository",
			"-m", "by https://github.com/g4s8/gitstrap").Run(); err != nil {
			return err
		}
		if err := exec.Command("git", "push", "origin", "master").Run(); err != nil {
			return err
		}
	}
	return nil
}

type templateContext struct {
	Repo     *github.Repository
	Gitstrap interface{}
}

func prompt(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", msg)
	text, _ := reader.ReadString('\n')
	a := strings.TrimSuffix(text, "\n")
	return strings.EqualFold(a, "y") || strings.EqualFold(a, "yes")
}
