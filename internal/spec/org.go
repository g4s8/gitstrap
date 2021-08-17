package spec

import (
	"github.com/google/go-github/v38/github"
)

type Org struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Company     string `yaml:"company,omitempty"`
	Blog        string `yaml:"blog,omitempty"`
	Location    string `yaml:"location,omitempty"`
	Email       string `yaml:"email,omitempty"`
	Twitter     string `yaml:"twitter,omitempty"`
	Verified    bool   `yaml:"verified,omitempty"`
}

func (o *Org) FromGithub(g *github.Organization) {
	o.Name = g.GetName()
	o.Company = g.GetCompany()
	o.Blog = g.GetBlog()
	o.Location = g.GetLocation()
	o.Email = g.GetEmail()
	o.Twitter = g.GetTwitterUsername()
	o.Description = g.GetDescription()
	o.Verified = g.GetIsVerified()
}

func (o *Org) ToGithub(g *github.Organization) error {
	g.Name = &o.Name
	g.Description = &o.Description
	g.Company = &o.Company
	g.Blog = &o.Blog
	g.Location = &o.Location
	g.Email = &o.Email
	g.TwitterUsername = &o.Twitter
	return nil
}
