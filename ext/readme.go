package ext

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
	"github.com/google/go-github/github"
	"log"
)

type badge struct {
	alt  string
	img  string
	link string
}

func (b *badge) markdown() string {
	return fmt.Sprintf("[![%s](%s)](%s)", b.alt, b.img, b.link)
}

func (b *badge) byName(n string, ctx *context.Context,
	p map[string]interface{}) error {
	owner := ctx.Owner()
	name := ctx.Name()
	switch n {
	case "rultor":
		b.alt = "DevOps By Rultor"
		b.img = fmt.Sprint("http://www.rultor.com/b/%s/%s", owner, name)
		b.link = fmt.Sprint("http://www.rultor.com/p/%s/%s", owner, name)
		return nil
	case "0crat":
		var proj string
		if err := tryStr(p, "project", &proj); err != nil {
			return err
		}
		b.alt = "Managed by Zerocracy"
		b.img = fmt.Sprintf("https://www.0crat.com/badge/%s.svg", proj)
		b.link = fmt.Sprintf("https://www.0crat.com/p/%s", proj)
		return nil
	case "releases":
		b.alt = "GitHub releases"
		b.img = fmt.Sprintf("https://img.shields.io/github/release/%s/%s.svg?label=version", owner, name)
		b.link = fmt.Sprintf("https://github.com/%s/%s/releases/latest)", owner, name)
		return nil
	default:
		return fmt.Errorf("Named extension '%s' was not found", name)
	}
}

func (b *badge) parse(p map[string]interface{}, ctx *context.Context) error {
	if name, has := p["name"]; has {
		var np map[string]interface{}
		if v, has := p["params"]; has {
			np = v.(map[string]interface{})
		} else {
			np = make(map[string]interface{})
		}
		nstr := name.(string)
		return b.byName(nstr, ctx, np)
	}
	if err := tryStr(p, "alt", &b.alt); err != nil {
		return fmt.Errorf("badge.alt: %s", err)
	}
	if err := tryStr(p, "img", &b.img); err != nil {
		return fmt.Errorf("badge.img: %s", err)
	}
	if err := tryStr(p, "link", &b.link); err != nil {
		return fmt.Errorf("badge.link: %s", err)
	}
	return nil
}

func tryStr(p map[string]interface{}, key string, out *string) error {
	val, has := p[key]
	if !has {
		return fmt.Errorf("no such key %s in %s", key, p)
	}
	s, ok := val.(string)
	if !ok {
		return fmt.Errorf("expected key %s to be string, got %s value", key, val)
	}
	*out = s
	return nil
}

type readme struct {
	title  string
	header string
	badges []badge
}

func (r *readme) parse(ctx *context.Context, params map[string]interface{}) error {
	tryStr(params, "title", &r.title)
	tryStr(params, "header", &r.header)
	if bmap, has := params["badges"]; has {
		lst := bmap.([]interface{})
		tmp := make([]map[interface{}]interface{}, len(lst), len(lst))
		for pos, item := range lst {
			tmp[pos] = item.(map[interface{}]interface{})
		}
		bs := make([]map[string]interface{}, len(lst), len(lst))
		for pos, mp := range tmp {
			bs[pos] = make(map[string]interface{})
			for k, v := range mp {
				bs[pos][k.(string)] = v
			}
		}
		// if bs, ok := bmap.([]map[string]interface{}); ok {
		r.badges = make([]badge, len(bs))
		for i, b := range bs {
			if err := r.badges[i].parse(b, ctx); err != nil {
				return err
			}
		}
		// } else {
		// 	return fmt.Errorf("incorrect badges type: %T %v", bmap, bmap)
		// }
	} else {
		r.badges = make([]badge, 0)
	}
	return nil
}

func (r *readme) content() string {
	var c string
	if r.title != "" {
		c += fmt.Sprintf("# %s\n\n", r.title)
	}
	cnt := 0
	for _, b := range r.badges {
		c += b.markdown()
		c += "\n"
		cnt++
		if cnt == 3 {
			c += "\n"
			cnt = 0
		}
	}
	if len(r.badges) != 0 && cnt != 0 && r.header != "" {
		c += "\n"
	}
	if r.header != "" {
		c += r.header
	}
	return c
}

func (r *readme) apply(ctx *context.Context) error {
	_, rsp, _ := ctx.Client.Repositories.GetReadme(ctx.Sync, ctx.Owner(), ctx.Name(), nil)
	exist := rsp.StatusCode != 404
	if exist {
		// readme file already exist, skipping this extension
		return nil
	}
	opts := new(github.RepositoryContentFileOptions)
	opts.Message = new(string)
	*opts.Message = "#1 - created readme file [gitstrap]\n\nSee: https://github.com/g4s8/gitstrap"
	opts.Content = []byte(r.content())
	res, _, err := ctx.Client.Repositories.CreateFile(ctx.Sync, ctx.Owner(), ctx.Name(), "README.md", opts)
	if err != nil {
		return fmt.Errorf("failed to upload readme: %s", err)
	}
	log.Printf("created readme file with at %s", res.Commit.GetSHA())
	return nil
}
