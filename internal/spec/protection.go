package spec

import (
	"github.com/google/go-github/v33/github"
)

// Protection rule of repositry branch
type Protection struct {
	// Require status checks for merge
	Require []string `yaml:"require,omitempty"`
	// Strict update with target branch is requried
	Strict bool `yaml:"strictUpdate"`
	// EnforceAdmins the same rules
	EnforceAdmins bool `yaml:"enforceAdmins,omitempty"`
	// LinearHistory is required for merging branch
	LinearHistory bool `yaml:"linearHistory,omitempty"`
	// ForcePush is allowed
	ForcePush bool `yaml:"forcePush,omitempty"`
	// CanDelete target branch
	CanDelete bool `yaml:"canDelete,omitempty"`
	// Permissions
	Permissions struct {
		// Restrict permissions is enabled
		Restrict bool `yaml:"restrict,omitempty"`
		// Users with push access
		Users []string `yaml:"users,omitempty"`
		// Teams with push access
		Teams []string `yaml:"teams,omitempty"`
		// Apps with push access
		Apps []string `yaml:"apps,omitempty"`
	} `yaml:"permissions,omitempty"`
}

func (bp *Protection) FromGithub(g *github.Protection) error {
	if c := g.RequiredStatusChecks; c != nil {
		bp.Require = make([]string, len(c.Contexts))
		for i, name := range c.Contexts {
			bp.Require[i] = name
		}
		bp.Strict = c.Strict
	}
	if e := g.EnforceAdmins; e != nil {
		bp.EnforceAdmins = e.Enabled
	}
	if l := g.RequireLinearHistory; l != nil {
		bp.LinearHistory = l.Enabled
	}
	if f := g.AllowForcePushes; f != nil {
		bp.ForcePush = f.Enabled
	}
	if d := g.AllowDeletions; d != nil {
		bp.CanDelete = d.Enabled
	}
	if r := g.Restrictions; r != nil {
		bp.Permissions.Restrict = true
		bp.Permissions.Users = make([]string, len(r.Users))
		for i, u := range r.Users {
			bp.Permissions.Users[i] = u.GetLogin()
		}
		bp.Permissions.Teams = make([]string, len(r.Teams))
		for i, t := range r.Teams {
			bp.Permissions.Teams[i] = t.GetSlug()
		}
		bp.Permissions.Apps = make([]string, len(r.Apps))
		for i, a := range r.Apps {
			bp.Permissions.Apps[i] = a.GetSlug()
		}
	}
	return nil
}
