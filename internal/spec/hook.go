package spec

import (
	"github.com/google/go-github/v33/github"
	"strconv"
)

type Hook struct {
	URL         string   `yaml:"url"`
	ContentType string   `yaml:"contentType"`
	InsecureSsl bool     `yaml:"insecureSsl,omitempty"`
	Secret      string   `yaml:"secret,omitempty"`
	Events      []string `yaml:"events,omitempty"`
	Active      bool     `yaml:"active"`
	Selector    struct {
		Repository   string `yaml:"repository,omitempty"`
		Organization string `yaml:"organization,omitempty"`
	} `yaml:"selector"`
}

const (
	hookCfgUrl         = "url"
	hookCfgContentType = "content_type"
	hookCfgInsecureSSL = "insecure_ssl"
)

func (h *Hook) FromGithub(g *github.Hook) error {
	if url, has := g.Config[hookCfgUrl]; has {
		h.URL = url.(string)
	}
	if ct, has := g.Config[hookCfgContentType]; has {
		h.ContentType = ct.(string)
	}
	if issl, has := g.Config[hookCfgInsecureSSL]; has {
		var err error
		h.InsecureSsl, err = strconv.ParseBool(issl.(string))
		if err != nil {
			return err
		}
	}
	h.Events = g.Events
	h.Active = g.GetActive()
	return nil
}

func (h *Hook) ToGithub(g *github.Hook) error {
	g.Config = make(map[string]interface{})
	g.Config[hookCfgUrl] = h.URL
	g.Config[hookCfgContentType] = h.ContentType
	g.Config[hookCfgInsecureSSL] = strconv.FormatBool(h.InsecureSsl)
	g.Events = h.Events
	g.Active = &h.Active
	return nil
}
