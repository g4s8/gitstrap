package spec

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v33/github"
)

// Metadata for spec
type Metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Repo        string            `yaml:"repo,omitempty"`
	Owner       string            `yaml:"owner,omitempty"`
	ID          *int64            `yaml:"id,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

func (m *Metadata) FromGithubRepo(r *github.Repository) {
	m.ID = r.ID
	m.Name = r.GetName()
}

func (m *Metadata) FromGithubOrg(o *github.Organization) {
	m.ID = o.ID
	m.Name = o.GetName()
}

func (m *Metadata) FromGithubTeam(t *github.Team) {
	m.ID = t.ID
	m.Name = *t.Slug
	m.Owner = *t.Organization.Login
}

func (m *Metadata) Info() string {
	sb := new(strings.Builder)
	if m.ID != nil {
		fmt.Fprintf(sb, "ID=%d ", *m.ID)
	}
	if m.Owner != "" {
		fmt.Fprintf(sb, "%s/%s ", m.Owner, m.Name)
	}
	if m.Name != "" {
		fmt.Fprintf(sb, "%s", m.Name)
	}
	return sb.String()
}
