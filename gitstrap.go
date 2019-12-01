package gitstrap

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/g4s8/gopwd"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Options - gitstrap options
type Options map[string]string

// Gitstrap - bootstrap tool
type Gitstrap interface {
	Run(opt Options) error
}

type strapCtx struct {
	cfg   *Config
	ctx   context.Context
	cli   *github.Client
	debug bool
}

func (ctx *strapCtx) String() string {
	return fmt.Sprintf("debug=%t, cfg=%+v", ctx.debug, ctx.cfg)
}

type strapCreate struct {
	base *strapCtx
}

func (s *strapCreate) String() string {
	return fmt.Sprintf("create={ctx={%s}}", s.base)
}

type strapDestr struct {
	base *strapCtx
}

type logTransport struct {
	origin http.RoundTripper
	tag    string
}

// @todo #37:30min Continue rafactoring gitstrap.go file
//  Move away all github logic, git functions and transport
//  structures. Then check again if refactoring is needed.

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("[%s] >>> %s %s", t.tag, req.Method, req.URL)
	if req.Body != nil {
		defer req.Body.Close()
		if data, err := ioutil.ReadAll(req.Body); err == nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			log.Print(string(data))
		}
	}
	rsp, err := t.origin.RoundTrip(req)
	if err != nil {
		log.Printf("[%s] %s ERR: %s", t.tag, req.URL, err)
	} else {
		log.Printf("[%s] %s <<< %d", t.tag, req.URL, rsp.StatusCode)
		if rsp.Body != nil {
			defer rsp.Body.Close()
			if data, err := ioutil.ReadAll(rsp.Body); err == nil {
				rsp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				log.Print(string(data))
			}
		}
	}
	return rsp, err
}

// New - make a gitstrap
func New(token string, action string, cfg *Config,
	debug bool) (Gitstrap, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	if debug {
		log.Printf("token: ***%s", token[:3])
	}
	tc := oauth2.NewClient(ctx, ts)
	if debug {
		tr := new(logTransport)
		tr.tag = "GH"
		if tc.Transport != nil {
			tr.origin = tc.Transport
		} else {
			tr.origin = http.DefaultTransport
		}
		tc.Transport = tr
	}
	cli := github.NewClient(tc)
	strap := &strapCtx{
		ctx:   ctx,
		cli:   cli,
		cfg:   cfg,
		debug: debug,
	}
	switch action {
	case "create":
		return &strapCreate{strap}, nil
	case "destroy":
		return &strapDestr{strap}, nil
	default:
		return nil, fmt.Errorf("unsupported action: '%s'", action)
	}
}

func (strap *strapCreate) Run(opt Options) error {
	name, err := strap.base.repoName()
	if err != nil {
		return err
	}
	repo := &github.Repository{
		Name:        &name,
		Description: strap.base.cfg.Gitstrap.Github.Repo.Description,
		Private:     strap.base.cfg.Gitstrap.Github.Repo.Private,
	}
	owner, err := getOwner(strap.base, opt)
	if strap.base.debug {
		fmt.Printf("RUN: as %s options = %s\n", owner, opt)
	}
	if err != nil {
		return ErrCompose(err, "failed to get current user")
	}
	repo, err = strap.base.createRepo(owner, repo, opt)
	if err != nil {
		return ErrCompose(err, "failed to create github repo")
	}
	if err := gitSync(repo); err != nil {
		return ErrCompose(err, "failed to sync git repo")
	}
	if err := strap.base.applyTemplates(repo); err != nil {
		return ErrCompose(err, "failed to apply templates")
	}
	if err := gitPush(repo); err != nil {
		return ErrCompose(err, "failed to push to remote")
	}
	if err := strap.base.addHooks(owner, repo); err != nil {
		return ErrCompose(err, "failed to add web-hooks")
	}
	if err := strap.base.addCollaborators(owner, repo); err != nil {
		return ErrCompose(err, "failed to add collaborators")
	}

	fmt.Printf("Created: https://github.com/%s/%s\n", owner, *repo.Name)

	return nil
}

func (s *strapCtx) createRepo(owner string, repo *github.Repository,
	opt Options) (*github.Repository, error) {
	fmt.Printf("Looking up for repo %s/%s\n", owner, *repo.Name)
	r, resp, _ := s.cli.Repositories.Get(s.ctx, owner, *repo.Name)
	exists := resp.StatusCode == 200
	_, accept := opt["accept"]
	if !exists && (accept || prompt("repository doesn't exist. Create?")) {
		if s.debug {
			fmt.Printf("creating repo: %s/%s\n\t%+v\n", owner, *repo.Name, repo)
		}
		org := ""
		if _, hasOrg := opt["org"]; hasOrg {
			org = owner
		}
		r, _, err := s.cli.Repositories.Create(s.ctx, org, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to create repo: %s", err)
		}
		fmt.Printf("Github repository %s has been created\n", *repo.Name)
		return r, nil
	} else if exists {
		fmt.Println("found")
	}
	return r, nil
}

func (strap *strapCtx) applyTemplates(repo *github.Repository) error {
	// apply templates
	tctx := &templateContext{repo, &strap.cfg.Gitstrap}
	for _, t := range strap.cfg.Gitstrap.Templates {
		tpl := template.New(t.Name)
		var data []byte
		var err error
		if t.Location != "" {
			data, err = readTemplate(t.Location)
			if err != nil {
				return err
			}
		} else if t.URL != "" {
			data, err = downloadTemplate(t.URL)
			if err != nil {
				return err
			}
		}
		if _, err = tpl.Parse(string(data)); err != nil {
			return fmt.Errorf("failed to parse template %s: %s", tpl.Name(), err)
		}
		fout, err := os.Create(t.Name)
		if err != nil {
			return fmt.Errorf("failed to open output file for template %s: %s", tpl.Name(), err)
		}
		if err = tpl.Execute(fout, tctx); err != nil {
			return fmt.Errorf("failed to execute template %s: %s", tpl.Name(), err)
		}
		fmt.Printf("Template %s applied\n", tpl.Name())
	}
	return nil
}

func readTemplate(name string) ([]byte, error) {
	tf, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open template file %s: %s", name, err)
	}
	data, err := ioutil.ReadAll(bufio.NewReader(tf))
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %s", name, err)
	}
	if err = tf.Close(); err != nil {
		return nil, fmt.Errorf("failed to close template file %s: %s", name, err)
	}
	return data, nil
}

func downloadTemplate(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download template from %s: %s", url, err)
	}
	data, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to read template body from %s: %s", url, err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close connection from %s: %s", url, err)
	}
	return data, nil
}

func (s *strapCtx) repoName() (string, error) {
	var name string
	if s.cfg.Gitstrap.Github.Repo.Name != nil {
		name = *s.cfg.Gitstrap.Github.Repo.Name
	} else {
		name, err := gopwd.Name()
		if err != nil {
			return name, errors.New("Failed to get PWD")
		}
	}
	return name, nil
}

func (strap *strapDestr) Run(opt Options) error {
	_, accept := opt["accept"]
	if !accept && !prompt("you are going to remove Github repository and local git repository. Are you sure?") {
		return nil
	}
	owner, err := getOwner(strap.base, opt)
	if err != nil {
		return ErrCompose(err, "failed to get current user")
	}
	name, err := strap.base.repoName()
	if err != nil {
		return err
	}
	fmt.Printf("Looking up for repo %s/%s... ", owner, name)
	_, resp, _ := strap.base.cli.Repositories.Get(strap.base.ctx, owner, name)
	exists := resp.StatusCode == 200
	if !exists {
		fmt.Printf("repository %s/%s not found\n", owner, name)
		os.Exit(1)
	}
	if _, err = strap.base.cli.Repositories.Delete(strap.base.ctx, owner, name); err != nil {
		return ErrCompose(err, "failed to delete repository")
	}
	fmt.Printf("Github repository %s/%s has been deleted\n", owner, name)
	if err = os.RemoveAll(".git"); err != nil {
		return ErrCompose(err, "Failed to remove git directory")
	}
	fmt.Println("Local git repository has been deleted")

	fmt.Println("Destroy: done")

	return nil
}

func (strap *strapCtx) addHooks(owner string, repo *github.Repository) error {
	for _, h := range strap.cfg.Gitstrap.Github.Repo.Hooks {
		ghkook := &github.Hook{
			URL:    &h.URL,
			Active: h.Active,
			Events: h.Events,
		}
		ghkook.Config = make(map[string]interface{})
		ghkook.Config["url"] = h.URL
		ghkook.Config["content_type"] = h.Type
		if _, _, err := strap.cli.Repositories.CreateHook(strap.ctx, owner, *repo.Name, ghkook); err != nil {
			return err
		}
		fmt.Printf("Webhook %s has been configured\n", h.URL)
	}
	return nil
}

func (strap *strapCtx) addCollaborators(owner string, repo *github.Repository) error {
	for _, clb := range strap.cfg.Gitstrap.Github.Repo.Collaborators {
		if _, err := strap.cli.Repositories.AddCollaborator(strap.ctx, owner, *repo.Name, clb, nil); err != nil {
			return err
		}
		fmt.Printf("Collaborator %s has been added\n", clb)
	}
	return nil
}

// getOwner returns current GitHub username,
// or organization name (if the -org flag was given)
func getOwner(strap *strapCtx, opt Options) (string, error) {
	org, hasOrg := opt["org"]
	if hasOrg {
		return org, nil
	}
	me, _, err := strap.cli.Users.Get(strap.ctx, "")
	if err != nil {
		return "", err
	}
	return *me.Login, nil
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
