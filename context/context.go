package context

import (
	"bytes"
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

// Options - gitstrap options
type Options map[string]string

// Context - gitstrap context
type Context struct {
	// Local repository path
	Path string
	// Github repository
	GhRepo *github.Repository
	// Sync context
	Sync context.Context
	// GitHub client
	Client *github.Client
	// Opt - gitstrap options
	Opt Options
	// internal
	owner *string
}

// New - create new context
func New(token string, path string, debug bool) *Context {
	ctx := context.Background()
	if debug {
		// print first chars of token if debug
		log.Printf("token: %s***", token[:3])
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
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
	cli := github.NewClient(tc)
	return &Context{
		Path:   path,
		Sync:   ctx,
		Client: cli,
		Opt:    Options(make(map[string]string)),
	}
}

type logTransport struct {
	origin http.RoundTripper
	tag    string
}

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("[%s] >>> %s %s", t.tag, req.Method, req.URL)
	if req.Body != nil {
		defer req.Body.Close()
		if data, err := ioutil.ReadAll(req.Body); err == nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			log.Print(string(data))
		}
	}
	rsp, err := t.origin.RoundTrip(req)
	if err != nil {
		log.Printf("[%s] %s ERR: %s", t.tag, req.URL, err)
	} else {
		log.Printf("[%s] %s <<< %d", t.tag, req.URL, rsp.StatusCode)
		if rsp.Body != nil {
			defer rsp.Body.Close()
			if data, err := ioutil.ReadAll(rsp.Body); err == nil {
				rsp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				log.Print(string(data))
			}
		}
	}
	return rsp, err
}
