package ext

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
	"github.com/google/go-github/github"
	"log"
	"strings"
)

type zpdd struct {
	verbose bool
	minWord int
}

func (p *zpdd) config() string {
	var c string
	c += "--source=.\n"
	c += fmt.Sprintf("--rule min-words:%d\n", p.minWord)
	if p.verbose {
		c += "--verbose\n"
	}
	return c
}

func (p *zpdd) parse(ctx *context.Context, params map[string]interface{}) error {
	if v, has := params["verbose"]; has {
		p.verbose = v == "true"
	} else {
		p.verbose = false
	}
	p.minWord = 20
	return nil
}

func (p *zpdd) uploadConfig(ctx *context.Context) error {
	_, _, rsp, _ := ctx.Client.Repositories.GetContents(ctx.Sync, ctx.Owner(), ctx.Name(), ".pdd", nil)
	if rsp.StatusCode != 404 {
		return nil
	}
	opts := new(github.RepositoryContentFileOptions)
	opts.Message = new(string)
	*opts.Message = "#1 - configured 0pdd\n\nSee: https://github.com/g4s8/gitstrap"
	opts.Content = []byte(p.config())
	res, _, err := ctx.Client.Repositories.CreateFile(ctx.Sync, ctx.Owner(), ctx.Name(), ".pdd", opts)
	if err != nil {
		return fmt.Errorf("failed to create .pdd config: %s", err)
	}
	log.Printf("created .pdd config at %s", res.GetSHA())
	return nil
}

func (p *zpdd) configureWebhook(ctx *context.Context) error {
	hks, _, err := ctx.Client.Repositories.ListHooks(ctx.Sync, ctx.Owner(), ctx.Name(), nil)
	if err != nil {
		return fmt.Errorf("failed to get repo hooks: %s", err)
	}
	exist := false
	for _, h := range hks {
		if strings.HasSuffix(h.GetURL(), "www.0pdd.com/hook/github") {
			exist = true
			break
		}
	}
	if exist {
		log.Print("0pdd web hook already exists")
		return nil
	}
	hook := new(github.Hook)
	hook.Config = make(map[string]interface{})
	hook.Config["url"] = "http://p.rehttp.net/http://www.0pdd.com/hook/github"
	hook.Config["content_type"] = "form"
	hook.Events = []string{"push"}
	if _, _, err := ctx.Client.Repositories.CreateHook(ctx.Sync, ctx.Owner(), ctx.Name(), hook); err != nil {
		return fmt.Errorf("failed to create 0pdd hook: %s", err)
	}
	log.Printf("0pdd webhook was configured sucessfully")
	return nil
}

func (p *zpdd) updateCollaborators(ctx *context.Context) error {
	users, _, err := ctx.Client.Repositories.ListCollaborators(ctx.Sync, ctx.Owner(), ctx.Name(), nil)
	if err != nil {
		return fmt.Errorf("failed to updated 0pdd collaborator: %s", err)
	}
	exist := false
	for _, u := range users {
		if u.GetLogin() == "0pdd" {
			exist = true
			break
		}
	}
	if exist {
		log.Print("0pdd user already exist")
		return nil
	}
	if _, err := ctx.Client.Repositories.AddCollaborator(ctx.Sync, ctx.Owner(), ctx.Name(), "0pdd", nil); err != nil {
		return fmt.Errorf("failed to add 0pdd collaborator: %s", err)
	}
	log.Print("0pdd collaborator has been invited")
	return nil
}

func (p *zpdd) apply(ctx *context.Context) error {
	if err := p.uploadConfig(ctx); err != nil {
		return err
	}
	if err := p.configureWebhook(ctx); err != nil {
		return err
	}
	if err := p.updateCollaborators(ctx); err != nil {
		return err
	}
	return nil
}
