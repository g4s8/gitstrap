package gitstrap

import (
	"context"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

// Gitstrap - main context
type Gitstrap struct {
	ctx   context.Context
	gh    *github.Client
	debug bool
	me    string
}

// New gitstrap context
func New(ctx context.Context, token string, debug bool) (*Gitstrap, error) {
	g := new(Gitstrap)
	g.debug = debug
	g.ctx = ctx
	if debug {
		// print first chars of token if debug
		log.Printf("Debug mode enabled: token='%s***%s'", token[:3], token[len(token)-2:])
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	if debug {
		// attach logging HTTP transport on debug
		tr := new(logTransport)
		tr.tag = "GH"
		if tc.Transport != nil {
			tr.origin = tc.Transport
		} else {
			tr.origin = http.DefaultTransport
		}
		tc.Transport = tr
	}
	g.gh = github.NewClient(tc)
	me, _, err := g.gh.Users.Get(g.ctx, "")
	if err != nil {
		return nil, err
	}
	g.me = me.GetLogin()
	return g, nil
}

func (g *Gitstrap) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(g.ctx, 25*time.Second)
}
