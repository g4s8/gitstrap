package main

import (
	"flag"
	"fmt"

	"github.com/g4s8/gitstrap"
	"os"
)

var (
	buildVersion string
	buildCommit  string
	buildDate    string
)

func main() {
	var config, token, org string
	var ver bool
	flag.StringVar(&token, "token", "", "Github API token")
	flag.StringVar(&config, "config", ".gitstrap.yaml", "Gitstrap config (default .gitstrap)")
	flag.StringVar(&org, "org", "", "Github organization (optional)")
	flag.BoolVar(&ver, "version", false, "Show version")
	flag.Parse()
	if ver {
		fmt.Printf("gitstrap version: %s\n"+
			"commit hash: %s\n"+
			"build date: %s\n", buildVersion, buildCommit, buildDate)
		os.Exit(0)
	}
	if token == "" {
		flagErr("Github token required")
	}

	if _, found := os.LookupEnv("SSH_AGENT_PID"); !found {
		fmt.Println("ssh-agent is not running. " +
			"You should start it before running gitstrap and add correct ssh key to be able to access Github repo via git")
		os.Exit(1)
	}
	cfg := &gitstrap.Config{}
	if err := cfg.ParseFile(config); err != nil {
		fatal(err)
	}
	action := flag.Arg(0)
	g, err := gitstrap.New(token, action, cfg)
	if err != nil {
		fatal(err)
	}
	options := gitstrap.Options(make(map[string]string))
	if org != "" {
		options["org"] = org
	}
	if err = g.Run(options); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	if _, xerr := fmt.Fprintf(os.Stderr, "%s\n", err); xerr != nil {
		fmt.Printf("Failed to print error: %s", xerr)
	}
	os.Exit(1)
}

func flagErr(msg string) {
	_, err := fmt.Fprintln(os.Stderr, msg)
	if err != nil {
		fmt.Printf("Failed to print errror: %s", err)
	}
	flag.Usage()
	os.Exit(1)
}
