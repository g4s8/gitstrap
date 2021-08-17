package spec

import (
	"github.com/google/go-github/v38/github"
)

// Protection rule of repositry branch
type Protection struct {
	// Checks represents required status checks for merge
	Checks []string `yaml:"checks,omitempty"`
	// Strict update with target branch is requried
	Strict bool `yaml:"strictUpdate,omitempty"`
	// Review represents pull request review enforcement
	Review struct {
		// Require pull request reviews enforcement of a protected branch.
		Require bool `yaml:"require,omitempty"`
		// Dismiss pull request review
		Dismiss struct {
			// Users who can dismiss review
			Users []string `yaml:"users,omitempty"`
			// Teams who can dismiss review
			Teams []string `yaml:"teams,omitempty"`
			// Automatically dismiss approving reviews when someone pushes a new commit.
			Stale bool `yaml:"stale,omitempty"`
		} `yaml:"dismiss,omitempty"`
		// RequireOwner blocks merging pull requests until code owners review them.
		RequireOwner bool `yaml:"requireOwner,omitempty"`
		// Count is the number of reviewers required to approve pull requests.
		Count int `yaml:"count,omitempty"`
	} `yaml:"review,omitempty"`
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
	// ConversationResolution, if set to true, requires all comments
	// on the pull request to be resolved before it can be merged to a protected branch.
	ConversationResolution bool `yaml:"conversationResolution,omitempty"`
}

func (bp *Protection) FromGithub(g *github.Protection) error {
	if c := g.RequiredStatusChecks; c != nil {
		bp.Checks = make([]string, len(c.Contexts))
		for i, name := range c.Contexts {
			bp.Checks[i] = name
		}
		bp.Strict = c.Strict
	}
	if p := g.RequiredPullRequestReviews; p != nil {
		bp.Review.Require = true
		if p.DismissStaleReviews {
			bp.Review.Dismiss.Stale = true
		}
		bp.Review.RequireOwner = p.RequireCodeOwnerReviews
		bp.Review.Count = p.RequiredApprovingReviewCount
		if r := p.DismissalRestrictions; r != nil {
			bp.Review.Dismiss.Users = make([]string, len(r.Users))
			for i, u := range r.Users {
				bp.Review.Dismiss.Users[i] = u.GetLogin()
			}
			bp.Review.Dismiss.Teams = make([]string, len(r.Teams))
			for i, t := range r.Teams {
				bp.Review.Dismiss.Teams[i] = t.GetSlug()
			}
		}
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
	if cr := g.RequiredConversationResolution; cr != nil {
		bp.ConversationResolution = cr.Enabled
	}
	return nil
}

func (bp *Protection) ToGithub(pr *github.ProtectionRequest) error {
	pr.EnforceAdmins = bp.EnforceAdmins
	pr.RequireLinearHistory = &bp.LinearHistory
	pr.AllowForcePushes = &bp.ForcePush
	pr.AllowDeletions = &bp.CanDelete
	if len(bp.Checks) != 0 || bp.Strict {
		pr.RequiredStatusChecks = bp.requiredChecksToGithub()
	}
	if bp.Review.Require {
		pr.RequiredPullRequestReviews = bp.reviewToGithub()
	}
	if bp.Permissions.Restrict {
		pr.Restrictions = bp.permissionsToGithub()
	}
	pr.RequiredConversationResolution = &bp.ConversationResolution
	return nil
}

func (bp *Protection) requiredChecksToGithub() *github.RequiredStatusChecks {
	c := new(github.RequiredStatusChecks)
	c.Contexts = *getEmptyIfNil(bp.Checks)
	c.Strict = bp.Strict
	return c
}

func (bp *Protection) reviewToGithub() *github.PullRequestReviewsEnforcementRequest {
	e := new(github.PullRequestReviewsEnforcementRequest)
	e.DismissalRestrictionsRequest = new(github.DismissalRestrictionsRequest)
	if d := bp.Review.Dismiss; d.Teams != nil || d.Users != nil {
		e.DismissalRestrictionsRequest.Teams = getEmptyIfNil(d.Teams)
		e.DismissalRestrictionsRequest.Users = getEmptyIfNil(d.Users)
	}
	e.DismissStaleReviews = bp.Review.Dismiss.Stale
	e.RequireCodeOwnerReviews = bp.Review.RequireOwner
	e.RequiredApprovingReviewCount = bp.Review.Count
	return e
}

func (bp *Protection) permissionsToGithub() *github.BranchRestrictionsRequest {
	r := new(github.BranchRestrictionsRequest)
	r.Teams = *getEmptyIfNil(bp.Permissions.Teams)
	r.Users = *getEmptyIfNil(bp.Permissions.Users)
	r.Apps = bp.Permissions.Apps
	return r
}

func getEmptyIfNil(slice []string) *[]string {
	if slice != nil {
		return &slice
	}
	return &[]string{}
}
