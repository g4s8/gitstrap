package spec

import (
	"github.com/google/go-github/v33/github"
)

type Team struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Slug        string `yaml:"slug,omitempty"`
	Permission  string `yaml:"permission,omitempty"`
	Privacy     string `yaml:"privacy,omitempty"`
}

func (t *Team) FromGithub(g *github.Team) error {
	t.Name = g.GetName()
	t.Description = g.GetDescription()
	t.Slug = g.GetSlug()
	t.Permission = g.GetPermission()
	t.Privacy = g.GetPrivacy()
	return nil
}

func (t *Team) ToGithub(g *github.NewTeam) error {
	g.Name = t.Name
	g.Description = &t.Description
	g.Privacy = &t.Privacy
	return nil
}
