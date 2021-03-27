package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/g4s8/gitstrap/internal/gitstrap"
	"github.com/g4s8/gitstrap/internal/spec"
	"github.com/urfave/cli/v2"
)

type errAggregate struct {
	errs []error
}

func (e *errAggregate) push(err error) {
	e.errs = append(e.errs, err)
}

func (e *errAggregate) ifAny() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e
}

func (e *errAggregate) Error() string {
	sb := new(strings.Builder)
	sb.WriteString("There are multiple errors occured:\n")
	for pos, err := range e.errs {
		sb.WriteString(fmt.Sprintf("  %d - %s\n", pos, err))
	}
	return sb.String()
}

type gtask func(g *gitstrap.Gitstrap, m *spec.Model) error

func cmdForEachModel(task gtask) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		token, err := resolveToken(ctx)
		if err != nil {
			return err
		}

		models, err := spec.ReadFile(ctx.String("file"))
		if err != nil {
			return err
		}
		debug := os.Getenv("DEBUG") != ""
		g, err := gitstrap.New(ctx.Context, token, debug)
		if err != nil {
			return err
		}
		if ctx.Bool("force") {
			for _, m := range models {
				m.Metadata.Annotations["force"] = "true"
			}
		}
		errs := new(errAggregate)
		for _, m := range models {
			if err := task(g, m); err != nil {
				errs.push(err)
			}
		}
		return errs.ifAny()
	}
}
