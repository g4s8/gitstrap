package spec

import (
	"github.com/google/go-github/v33/github"
)

type Hook struct {
	ID          *int64   `yaml:"id,omitempty"`
	URL         string   `yaml:"url"`
	ContentType string   `yaml:"contentType"`
	InsecureSsl bool     `yaml:"insecureSsl,omitempty"`
	Secret      string   `yaml:"secret,omitempty"`
	Events      []string `yaml:"events,omitempty"`
	Active      bool     `yaml:"active,omitempty"`
}

func (h *Hook) FromGithub(g *github.Hook) {
	h.ID = g.ID
	if url, has := g.Config["url"]; has {
		h.URL = url.(string)
	}
	if ct, has := g.Config["content_type"]; has {
		h.ContentType = ct.(string)
	}
	if issl, has := g.Config["insecure_ssl"]; has {
		h.InsecureSsl = issl.(string) == "true"
	}
	if sec, has := g.Config["secret"]; has {
		h.Secret = sec.(string)
	}
	h.Events = g.Events
	h.Active = g.GetActive()
}
