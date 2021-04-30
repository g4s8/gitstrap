package gitstrap

import (
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
)

func (g *Gitstrap) getOwner(m *spec.Model) string {
	owner := m.Metadata.Owner
	if owner == "" {
		owner = g.me
	}
	return owner
}

//resolveOrg determines whether the owner is an organization.
//If the owner is the same as the authorized user, it returns an empty string.
//Otherwise it returns owner, because githab only allows repositories
//to be created in a personal account or in an organization the user is a member of.
func (g *Gitstrap) resolveOrg(m *spec.Model) string {
	if m.Metadata.Owner == g.me {
		return ""
	}
	return m.Metadata.Owner
}

type errNotSpecified struct {
	field string
}

func (e *errNotSpecified) Error() string {
	return fmt.Sprintf("%v is not specified", e.field)
}

func getSpecifiedOwner(m *spec.Model) (string, error) {
	owner := m.Metadata.Owner
	if owner == "" {
		return "", &errNotSpecified{"Owner"}
	}
	return owner, nil
}

func getSpecifiedName(m *spec.Model) (string, error) {
	name := m.Metadata.Name
	if name == "" {
		return "", &errNotSpecified{"Name"}
	}
	return name, nil
}

func getSpecifiedID(m *spec.Model) (*int64, error) {
	ID := m.Metadata.ID
	if ID == nil {
		return nil, &errNotSpecified{"ID"}
	}
	return ID, nil
}
