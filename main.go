package main

import (
	"./cfg"
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var (
	token string
	conf  *cfg.YamlConfig
	ctx   context.Context
	cli   *github.Client
)

func main() {
	var config string
	flag.StringVar(&token, "token", "", "Github API token")
	flag.StringVar(&config, "config", ".gitstrap.yaml", "Gitstrap config (default .gitstrap)")
	flag.Parse()
	if token == "" {
		flagErr("Github token required")
	}

	if _, found := os.LookupEnv("SSH_AGENT_PID"); !found {
		fmt.Println("ssh-agent is not running. " +
			"You should start it before gitstrap and add correct ssh key to be able access Github repo via git")
		os.Exit(1)
	}

	fconf, err := os.Open(config)
	fatal("Failed to open config file", err)
	conf = &cfg.YamlConfig{}
	err = conf.Parse(bufio.NewReader(fconf))
	fatal("Failed to read config", err)
	err = fconf.Close()
	fatal("Failed to close config file", err)

	ctx = context.Background()
	cli = githubCli(ctx, token)

	switch cmd := flag.Arg(0); cmd {
	case "create":
		create()
	case "destroy":
		destroy()
	default:
		fmt.Printf("unknown command %s - use create or destroy")
		os.Exit(1)
	}
}

func create() {
	var repo github.Repository
	repo.Name = conf.Gitstrap.Github.Repo.Name
	repo.Description = conf.Gitstrap.Github.Repo.Description
	repo.Private = conf.Gitstrap.Github.Repo.Private
	repo.AutoInit = conf.Gitstrap.Github.Repo.AutoInit

	// find current user
	me, _, err := cli.Users.Get(ctx, "")
	fatal("Failed to get current user", err)

	// find or create repo
	fmt.Printf("Looking up for repo %s/%s... ", *me.Login, *repo.Name)
	r, resp, _ := cli.Repositories.Get(ctx, *me.Login, *repo.Name)
	exists := resp.StatusCode == 200
	if !exists && prompt("repository doesn't exist. Create?") {
		r, _, err := cli.Repositories.Create(ctx, "", &repo)
		fatal(fmt.Sprintf("Failed to create repo %s", *repo.Name), err)
		repo = *r
		fmt.Printf("Github repository %s has been created\n", *repo.Name)
	} else if exists {
		fmt.Println("found")
		repo = *r
	}

	gitSync(&repo)

	// apply templates
	tctx := &templateContext{&repo, &conf.Gitstrap}
	for _, t := range conf.Gitstrap.Templates {
		tpl := template.New(t.Name)
		tf, err := os.Open(t.Location)
		fatal(fmt.Sprintf("Failed to open template file %s", t.Location), err)
		data, err := ioutil.ReadAll(bufio.NewReader(tf))
		fatal(fmt.Sprintf("Failed to read template file %s", t.Location), err)
		err = tf.Close()
		fatal("Failed to close template file", err)
		_, err = tpl.Parse(string(data))
		fatal(fmt.Sprintf("Failed to parse template %s", tpl.Name()), err)
		fout, err := os.Create(t.Name)
		fatal(fmt.Sprintf("Failed to open output file for template %s", tpl.Name()), err)
		err = tpl.Execute(fout, tctx)
		fatal(fmt.Sprintf("Failed to execute template %s", tpl.Name()), err)
		fmt.Printf("Template %s applied\n", tpl.Name())
	}

	gitPush(&repo)

	addHooks(me, &repo, &conf.Gitstrap)
	addCollaborators(me, &repo, &conf.Gitstrap)

	fmt.Println("Create: done")
}

func destroy() {
	if !prompt("you are going to remove Github repository and local git repository. Are you sure?") {
		return
	}
	me, _, err := cli.Users.Get(ctx, "")
	fatal("Failed to get current user", err)
	name := *conf.Gitstrap.Github.Repo.Name
	fmt.Printf("Looking up for repo %s/%s... ", *me.Login, name)
	_, resp, _ := cli.Repositories.Get(ctx, *me.Login, name)
	exists := resp.StatusCode == 200
	if !exists {
		fmt.Printf("repository %s/%s not found\n", *me.Login, name)
		os.Exit(1)
	}
	_, err = cli.Repositories.Delete(ctx, *me.Login, name)
	fatal("Failed to delete repository", err)
	fmt.Printf("Github repository %s/%s has been deleted\n", *me.Login, name)
	err = os.RemoveAll(".git")
	fatal("Failed to remove git directory", err)
	fmt.Println("Local git repository has been deleted")

	fmt.Println("Destroy: done")
}

func addHooks(me *github.User, repo *github.Repository, g *cfg.Gitstrap) {
	for _, h := range g.Github.Repo.Hooks {
		ghkook := &github.Hook{
			URL:    &h.Url,
			Active: h.Active,
			Events: h.Events,
		}
		ghkook.Config = make(map[string]interface{})
		ghkook.Config["url"] = h.Url
		ghkook.Config["content_type"] = h.Type
		_, _, err := cli.Repositories.CreateHook(ctx, *me.Login, *repo.Name, ghkook)
		fatal("Failed to create webhook", err)
		fmt.Printf("Webhook %s has been configured\n", h.Url)
	}
}

func addCollaborators(me *github.User, repo *github.Repository, g *cfg.Gitstrap) {
	for _, clb := range g.Github.Repo.Collaborators {
		_, err := cli.Repositories.AddCollaborator(ctx, *me.Login, *repo.Name, clb, nil)
		fatal("Failed to add collaborator", err)
		fmt.Printf("Collaborator %s has been added\n", clb)
	}
}

func gitSync(repo *github.Repository) {
	runCmd(exec.Command("git", "init", "."))
	runCmd(exec.Command("git", "remote", "add", "origin", *repo.SSHURL))
	runCmd(exec.Command("git", "fetch", "origin"))
	fmt.Printf("Github repository %s has been fetched\n", *repo.SSHURL)

	if prompt("Run `git pull` (if repository is not empty)?") {
		runCmd(exec.Command("git", "pull", "origin", "master"))
		fmt.Printf("pulled\n")
	}
}

func gitPush(repo *github.Repository) {
	if prompt("Templates has been applied. Do you want to commit & push?") {
		runCmd(exec.Command("git", "add", "."))
		runCmd(exec.Command("git", "commit", "-m", "[gitstrap] bootstrap repository"))
		runCmd(exec.Command("git", "push", "origin", "master"))
	}
}

func fatal(msg string, err error) {
	if err == nil {
		return
	}
	if _, xerr := fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err); xerr != nil {
		fmt.Printf("Failed to print error: %s", xerr)
	}
	os.Exit(1)
}

type templateContext struct {
	Repo     *github.Repository
	Gitstrap *cfg.Gitstrap
}

func prompt(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", msg)
	text, _ := reader.ReadString('\n')
	a := strings.TrimSuffix(text, "\n")
	return strings.EqualFold(a, "y") || strings.EqualFold(a, "yes")
}

func runCmd(c *exec.Cmd) {
	if err := c.Run(); err != nil {
		if _, xerr := fmt.Fprintf(os.Stderr, "command %s failed: %s\n", c.Path, err); xerr != nil {
			fmt.Printf("Failed to print error: %s", xerr)
		}
		os.Exit(1)
	}
}

func githubCli(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func flagErr(msg string) {
	_, err := fmt.Fprintln(os.Stderr, msg)
	if err != nil {
		fmt.Printf("Failed to print errror: %s", err)
	}
	flag.Usage()
	os.Exit(1)
}
