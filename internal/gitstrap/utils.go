package gitstrap

import "github.com/g4s8/gitstrap/internal/spec"

func (g *Gitstrap) getOwner(m *spec.Model) string {
	owner := m.Metadata.Owner
	if owner == "" {
		owner = g.me
	}
	return owner
}
