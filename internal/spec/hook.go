package spec

import (
	"strconv"

	"github.com/google/go-github/v38/github"
)

type Hook struct {
	URL         string   `yaml:"url" default:"http://example.com/hook"`
	ContentType string   `yaml:"contentType" default:"json"`
	InsecureSsl bool     `yaml:"insecureSsl,omitempty" default:"false"`
	Secret      string   `yaml:"secret,omitempty"`
	Events      []string `yaml:"events,omitempty" default:"[\"push\"]"`
	Active      bool     `yaml:"active" default:"true"`
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
