package spec

import (
	"github.com/google/go-github/v36/github"
)

const (
	FeatureIssues    = "issues"
	FeatureWiki      = "wiki"
	FeaturePages     = "pages"
	FeatureProjects  = "projects"
	FeatureDownloads = "downloads"
)

const (
	RepoVisibilityPublic  = "public"
	RepoVisibilityPrivate = "private"
)

const (
	MergeCommit = "commit"
	MergeRebase = "rebase"
	MergeSquash = "squash"
)

// Repo spec
type Repo struct {
	Description         *string  `yaml:"description,omitempty"`
	Homepage            *string  `yaml:"homepage,omitempty"`
	DefaultBranch       string   `yaml:"defaultBranch,omitempty" default:"master"`
	MergeStrategy       []string `yaml:"mergeStrategy,omitempty" default:"[\"merge\"]"`
	DeleteBranchOnMerge *bool    `yaml:"deleteBranchOnMerge,omitempty"`
	Topics              []string `yaml:"topics,omitempty"`
	Archived            *bool    `yaml:"archived,omitempty"`
	Disabled            *bool    `yaml:"disabled,omitempty"`
	License             *string  `yaml:"license,omitempty"`
	Visibiliy           *string  `yaml:"visibility,omitempty" default:"public"`
	Features            []string `yaml:"features,omitempty"`
}

func (spec *Repo) FromGithub(repo *github.Repository) {
	spec.Description = repo.Description
	spec.Homepage = repo.Homepage
	spec.DefaultBranch = repo.GetDefaultBranch()
	if l := repo.GetLicense(); l != nil {
		spec.License = l.Key
	}
	spec.MergeStrategy = make([]string, 0, 3)
	if repo.GetAllowMergeCommit() {
		spec.MergeStrategy = append(spec.MergeStrategy, "commit")
	}
	if repo.GetAllowRebaseMerge() {
		spec.MergeStrategy = append(spec.MergeStrategy, "rebase")
	}
	if repo.GetAllowSquashMerge() {
		spec.MergeStrategy = append(spec.MergeStrategy, "squash")
	}
	spec.DeleteBranchOnMerge = repo.DeleteBranchOnMerge
	spec.Topics = repo.Topics
	if repo.GetArchived() {
		spec.Archived = repo.Archived
	}
	if repo.GetDisabled() {
		spec.Disabled = repo.Disabled
	}
	spec.Visibiliy = new(string)
	if repo.GetPrivate() {
		*spec.Visibiliy = RepoVisibilityPrivate
	} else {
		*spec.Visibiliy = RepoVisibilityPublic
	}
	// issues, wiki, pages, projects, downloads
	spec.Features = make([]string, 0, 5)
	if repo.GetHasIssues() {
		spec.Features = append(spec.Features, FeatureIssues)
	}
	if repo.GetHasWiki() {
		spec.Features = append(spec.Features, FeatureWiki)
	}
	if repo.GetHasPages() {
		spec.Features = append(spec.Features, FeaturePages)
	}
	if repo.GetHasProjects() {
		spec.Features = append(spec.Features, FeatureProjects)
	}
	if repo.GetHasDownloads() {
		spec.Features = append(spec.Features, FeatureDownloads)
	}
}

func (s *Repo) ToGithub(r *github.Repository) error {
	r.Description = s.Description
	r.Homepage = s.Homepage
	r.DefaultBranch = &s.DefaultBranch
	r.AllowMergeCommit = new(bool)
	r.AllowRebaseMerge = new(bool)
	r.AllowSquashMerge = new(bool)
	for _, ms := range s.MergeStrategy {
		switch ms {
		case "commit":
			*r.AllowMergeCommit = true
		case "rebase":
			*r.AllowRebaseMerge = true
		case "squash":
			*r.AllowSquashMerge = true
		}
	}
	if !r.GetAllowMergeCommit() && !r.GetAllowRebaseMerge() && !r.GetAllowSquashMerge() {
		// enabled at least merge commit
		*r.AllowMergeCommit = true
	}
	r.DeleteBranchOnMerge = s.DeleteBranchOnMerge
	r.Topics = s.Topics
	r.Archived = s.Archived
	r.Disabled = s.Disabled
	if s.License != nil {
		r.License = new(github.License)
		r.License.Key = s.License
	}
	r.Private = new(bool)
	if s.Visibiliy == nil || *s.Visibiliy == "public" {
		*r.Private = false
	} else if s.Visibiliy != nil && *s.Visibiliy == "private" {
		*r.Private = true
	}
	r.HasIssues = new(bool)
	r.HasWiki = new(bool)
	r.HasPages = new(bool)
	r.HasProjects = new(bool)
	r.HasDownloads = new(bool)
	for _, f := range s.Features {
		// issues, wiki, pages, projects, downloads
		switch f {
		case "issues":
			*r.HasIssues = true
		case "wiki":
			*r.HasWiki = true
		case "pages":
			*r.HasPages = true
		case "projects":
			*r.HasProjects = true
		case "downloads":
			*r.HasDownloads = true
		}
	}
	return nil
}
