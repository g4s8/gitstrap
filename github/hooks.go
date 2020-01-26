package github

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
	"github.com/google/go-github/github"
	"log"
)

// RepoHooks - GitHub repository webhooks
type RepoHooks []RepoHook

// RepoHook - GitHub repository webhook
type RepoHook struct {
	URL    string   `yaml:"url"`
	Type   string   `yaml:"type"`
	Events []string `yaml:"events"`
	Active *bool    `yaml:"active"`
	// GitHub hook ID, internal use only
	gid *int64
}

// Apply - apply webhook configuration
func (hks RepoHooks) Apply(ctx *context.Context) error {
	own := repoOwner(ctx.GhRepo)
	name := *ctx.GhRepo.Name
	rhks, _, err := ctx.Client.Repositories.ListHooks(ctx.Sync, own, name, nil)
	if err != nil {
		return fmt.Errorf("failed to list repo hooks: %s", err)
	}

	// convert GitHub response to gitstrap hook data, map by URL to compare it later
	old := make(map[string]RepoHook, len(rhks))
	for _, h := range rhks {
		var hook RepoHook
		if err := hook.fromGithub(h); err != nil {
			return err
		}
		old[hook.URL] = hook
	}

	// nwe=new - map new hooks by URL to compare it
	nwe := make(map[string]RepoHook, len(hks))
	for _, h := range hks {
		nwe[h.URL] = h
	}

	// add - new hooks to add, update - existing hooks to update, delete - existng hooks to delete
	// unchanged - existing hooks not changed
	add := make([]RepoHook, 0)
	update := make([]RepoHook, 0)
	del := make([]RepoHook, 0)
	unchanged := make([]RepoHook, 0)

	for u, h := range nwe {
		if ex, has := old[u]; has {
			// if hook exists and differ from new hook
			if ex.cmp(&h) != 0 {
				// update existing hook (keep ex.gid)
				ex.update(&h)
				update = append(update, ex)
			} else {
				// if hook exists but same as new hook
				// save as unchanged
				unchanged = append(unchanged, ex)
			}
		} else {
			// if hook doesn't exist, add it
			add = append(add, h)
		}
	}
	for u, h := range old {
		if _, has := nwe[u]; !has {
			// if old hook doesn't exist in new hooks, remove it
			del = append(del, h)
		}
	}

	// removing hooks
	for _, h := range del {
		if _, err := ctx.Client.Repositories.DeleteHook(ctx.Sync, own, name, *h.gid); err != nil {
			return fmt.Errorf("failed to delete webhook %s (%s)", h.URL, err)
		}
		log.Printf("removed webhook %s", h.URL)
	}

	// updating hooks
	for _, h := range update {
		gh := new(github.Hook)
		h.toGithub(gh)
		if _, _, err := ctx.Client.Repositories.EditHook(ctx.Sync, own, name, *h.gid, gh); err != nil {
			return fmt.Errorf("failed to update webhook %s (%s)", h.URL, err)
		}
		log.Printf("updated webhook %s", h.URL)
	}

	// adding new hooks
	for _, h := range add {
		gh := new(github.Hook)
		h.toGithub(gh)
		if _, _, err := ctx.Client.Repositories.CreateHook(ctx.Sync, own, name, gh); err != nil {
			return fmt.Errorf("failed to create webhook %s (%s)", h.URL, err)
		}
		log.Printf("created webhook %s", h.URL)
	}

	// print unchanged
	for _, h := range unchanged {
		log.Printf("webhook unchanged %s", h.URL)
	}

	return nil
}

func (h *RepoHook) toGithub(gh *github.Hook) {
	gh.Config = make(map[string]interface{})
	gh.Config["url"] = h.URL
	gh.Config["content_type"] = h.Type
	gh.Active = new(bool)
	*gh.Active = h.IsActive()
	gh.Events = h.Events
	gh.ID = h.gid
}

// IsActive - true if webhook is active
func (h *RepoHook) IsActive() bool {
	return h.Active == nil || *h.Active
}

// SetActive - change active flag
func (h *RepoHook) SetActive(val bool) {
	*h.Active = val
}

func (h *RepoHook) update(val *RepoHook) {
	h.SetActive(val.IsActive())
	h.Events = val.Events
	h.Type = val.Type
}

func (l *RepoHook) cmp(r *RepoHook) int {
	var ch int
	if l.URL != r.URL {
		ch++
	}
	if l.IsActive() != r.IsActive() {
		ch++
	}
	if !sameStringSlice(l.Events, r.Events) {
		ch++
	}
	if l.Type != r.Type {
		ch++
	}
	return ch
}

// https://stackoverflow.com/a/36000696
func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}

func (h *RepoHook) fromGithub(gh *github.Hook) error {
	var url, tpe string
	if val, has := gh.Config["url"]; has {
		url = val.(string)
	}
	if val, has := gh.Config["content_type"]; has {
		tpe = val.(string)
	}
	if url == "" || tpe == "" {
		return fmt.Errorf("broken hook: no URL or type; "+
			"URL='%s' type='%s'", url, tpe)
	}
	h.URL = url
	h.Type = tpe
	h.Events = gh.Events
	h.Active = gh.Active
	h.gid = gh.ID
	return nil
}

func repoOwner(r *github.Repository) string {
	if r.Organization != nil {
		return *r.Organization.Login
	}
	return *r.Owner.Login
}

func (hks RepoHooks) contains(url string) bool {
	for _, h := range hks {
		if h.URL == url {
			return true
		}
	}
	return false
}
