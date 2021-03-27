package spec

import (
	"github.com/google/go-github/v33/github"
)

type Org struct {
	Company     string `yaml:"company,omitempty"`
	Blog        string `yaml:"blog,omitempty"`
	Location    string `yaml:"location,omitempty"`
	Email       string `yaml:"email,omitempty"`
	Twitter     string `yaml:"twitter,omitempty"`
	Description string `yaml:"description,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Verified    bool   `yaml:"verified,omitempty"`
}

func (o *Org) FromGithub(g *github.Organization) {
	o.Company = g.GetCompany()
	o.Blog = g.GetBlog()
	o.Location = g.GetLocation()
	o.Email = g.GetEmail()
	o.Twitter = g.GetTwitterUsername()
	o.Description = g.GetDescription()
	o.Type = g.GetType()
	o.Verified = g.GetIsVerified()
}
