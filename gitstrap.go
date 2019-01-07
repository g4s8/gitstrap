package gitstrap

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	// V1 - first version of config
	V1 = "v1"
)

// Config - gitstrap config
type Config struct {
	Gitstrap *struct {
		Version string `yaml:"version"`
		Github  *struct {
			Repo *struct {
				Name        *string `yaml:"name"`
				Description *string `yaml:"description"`
				Private     *bool   `yaml:"private"`
				AutoInit    *bool   `yaml:"autoInit"`
				Hooks       []struct {
					URL    string   `yaml:"url"`
					Type   string   `yaml:"type"`
					Events []string `yaml:"events"`
					Active *bool    `yaml:"active"`
				} `yaml:"hooks"`
				Collaborators []string `yaml:"collaborators"`
			} `yaml:"repo"`
		} `yaml:"github"`
		Templates []struct {
			Name     string `yaml:"name"`
			Location string `yaml:"location"`
			URL      string `yaml:"url"`
		} `yaml:"templates"`
		Params map[string]string `yaml:"params"`
	} `yaml:"gitstrap"`
}

// ParseReader - parse config from reader
func (y *Config) ParseReader(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(y)
	if y.Gitstrap.Version != V1 {
		return fmt.Errorf("Unsupported version: %s", y.Gitstrap.Version)
	}
	return err
}

// ParseFile - parse config from file
func (y *Config) ParseFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("Failed to open config file: %s", err)
	}
	if err = y.ParseReader(bufio.NewReader(f)); err != nil {
		return err
	}
	err = f.Close()
	return err
}

// Options - gitstrap options
type Options map[string]string

// Gitstrap - bootstrap tool
type Gitstrap interface {
	Run(opt Options) error
}

type strapCtx struct {
	cfg *Config
	ctx context.Context
	cli *github.Client
}

type strapCreate struct {
	base *strapCtx
}

type strapErr struct {
	strap string
	msg   string
	cause error
}

type strapDestr struct {
	base *strapCtx
}

func (err *strapErr) Error() string {
	return fmt.Sprintf("[%s] %s: %s", err.strap, err.msg, err.cause)
}

func (strap *strapCreate) err(msg string, cause error) error {
	return &strapErr{"create", msg, cause}
}

func (strap *strapDestr) err(msg string, cause error) error {
	return &strapErr{"destroy", msg, cause}
}

// New - make a gitstrap
func New(token string, action string, cfg *Config) (Gitstrap, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli := github.NewClient(tc)
	strap := &strapCtx{
		ctx: ctx,
		cli: cli,
		cfg: cfg,
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
	// @todo #5:30min Continue refactoring.
	//  Refactor strapCreate.Run and strapDestr.Run functions.
	//  Split create and destroy logic to independend steps
	//  (such as create github repo, init git repo, etc),
	//  Implement each actions as a chain of these steps.
	repo := &github.Repository{
		Name:        strap.base.cfg.Gitstrap.Github.Repo.Name,
		Description: strap.base.cfg.Gitstrap.Github.Repo.Description,
		Private:     strap.base.cfg.Gitstrap.Github.Repo.Private,
	}

	// find current user
	me, _, err := strap.base.cli.Users.Get(strap.base.ctx, "")
	if err != nil {
		return strap.err("failed to get current user", err)
	}
	repo, err = strap.base.createRepo(*me.Login, repo)
	if err != nil {
		return strap.err("failed to create github repo", err)
	}
	if err := gitSync(repo); err != nil {
		return strap.err("failed to sync git repo", err)
	}
	if err := strap.base.applyTemplates(repo); err != nil {
		return strap.err("failed to apply templates", err)
	}
	if err := gitPush(repo); err != nil {
		return strap.err("failed to push to remote", err)
	}
	if err := strap.base.addHooks(me, repo); err != nil {
		return strap.err("failed to add web-hooks", err)
	}
	if err := strap.base.addCollaborators(me, repo); err != nil {
		return strap.err("failed to add collaborators", err)
	}

	fmt.Println("Create: done")

	return nil
}

func (strap *strapCtx) createRepo(owner string, repo *github.Repository) (*github.Repository, error) {
	fmt.Printf("Looking up for repo %s/%s... ", owner, *repo.Name)
	r, resp, _ := strap.cli.Repositories.Get(strap.ctx, owner, *repo.Name)
	exists := resp.StatusCode == 200
	if !exists && prompt("repository doesn't exist. Create?") {
		r, _, err := strap.cli.Repositories.Create(strap.ctx, "", repo)
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

func (strap *strapDestr) Run(opt Options) error {
	if !prompt("you are going to remove Github repository and local git repository. Are you sure?") {
		return nil
	}
	me, _, err := strap.base.cli.Users.Get(strap.base.ctx, "")
	if err != nil {
		return strap.err("failed to get current user", err)
	}
	name := *strap.base.cfg.Gitstrap.Github.Repo.Name
	fmt.Printf("Looking up for repo %s/%s... ", *me.Login, name)
	_, resp, _ := strap.base.cli.Repositories.Get(strap.base.ctx, *me.Login, name)
	exists := resp.StatusCode == 200
	if !exists {
		fmt.Printf("repository %s/%s not found\n", *me.Login, name)
		os.Exit(1)
	}
	if _, err = strap.base.cli.Repositories.Delete(strap.base.ctx, *me.Login, name); err != nil {
		return strap.err("failed to delete repository", err)
	}
	fmt.Printf("Github repository %s/%s has been deleted\n", *me.Login, name)
	if err = os.RemoveAll(".git"); err != nil {
		return strap.err("Failed to remove git directory", err)
	}
	fmt.Println("Local git repository has been deleted")

	fmt.Println("Destroy: done")

	return nil
}

func (strap *strapCtx) addHooks(me *github.User, repo *github.Repository) error {
	for _, h := range strap.cfg.Gitstrap.Github.Repo.Hooks {
		ghkook := &github.Hook{
			URL:    &h.URL,
			Active: h.Active,
			Events: h.Events,
		}
		ghkook.Config = make(map[string]interface{})
		ghkook.Config["url"] = h.URL
		ghkook.Config["content_type"] = h.Type
		if _, _, err := strap.cli.Repositories.CreateHook(strap.ctx, *me.Login, *repo.Name, ghkook); err != nil {
			return err
		}
		fmt.Printf("Webhook %s has been configured\n", h.URL)
	}
	return nil
}

func (strap *strapCtx) addCollaborators(me *github.User, repo *github.Repository) error {
	for _, clb := range strap.cfg.Gitstrap.Github.Repo.Collaborators {
		if _, err := strap.cli.Repositories.AddCollaborator(strap.ctx, *me.Login, *repo.Name, clb, nil); err != nil {
			return err
		}
		fmt.Printf("Collaborator %s has been added\n", clb)
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
		if err := exec.Command("git", "commit", "-m", "[gitstrap] bootstrap repository").Run(); err != nil {
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
