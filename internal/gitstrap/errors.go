package gitstrap

import (
	"errors"
	"fmt"

	"github.com/g4s8/gitstrap/internal/spec"
)

var (
	errHookSelectorEmpty = errors.New("hook selector is empty: requires repository or organization")
	errHookIdRequired    = errors.New("Hook metadata ID required")
)

type errReadmeNotExists struct {
	owner, repo string
}

func (e *errReadmeNotExists) Error() string {
	return fmt.Sprintf("README.md `%s/%s` doesn't exist", e.owner, e.repo)
}

type errReadmeExists struct {
	owner, repo string
}

func (e *errReadmeExists) Error() string {
	return fmt.Sprintf("README.md already exists in %s/%s (try --force for replacing it)", e.owner, e.repo)
}

type errReadmeNotFile struct {
	rtype string
}

func (e *errReadmeNotFile) Error() string {
	return fmt.Sprintf("README is no a file: `%s`", e.rtype)
}

type errUnsupportModelKind struct {
	kind spec.Kind
}

func (e *errUnsupportModelKind) Error() string {
	return fmt.Sprintf("Unsupported model kind: `%s`", e.kind)
}

type errNotSpecified struct {
	field string
}

func (e *errNotSpecified) Error() string {
	return fmt.Sprintf("%v is not specified", e.field)
}
