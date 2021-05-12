
   *The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
   "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
   document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).*

Gitsrap works mostly with resource specification yaml documents. Each file should contain at least one document and may contain multiple documents separated by `---` yaml document separator.

Each specification document is a yaml mapping, it consists of these keys:
 - `version` - specification document version
 - `kind` - specification kind
 - `metadata` - specification metadata
 - `spec` - resource specification

Version should be `v2`. Kind is a string which describes the specification kind, it must be one of:
 - [Repository](#Repository)
 - [Organization](#Organization)
 - [WebHook](#WebHook)
 - [Readme](#Readme)
 - [Team](#Team)

A `spec` may have `selector` key to attach to the correct resource.

Metadata is differ for specification, but most of them has these keys:
 - `name` - resource name (e.g. repository name or organization name)
 - `id` - resource ID (e.g. repository, webhook, organization ID)
 - `owner` - resource owner, e.g. repository owner

## Repository

Describes GitHub repository resource, it has:
 - `description` (string) - repository description
 - `homepage` (string) - repository home page URI
 - `defaultBranch` (string, default: `"master"`) - the name of the default branch
 - `mergeStrategy` (list of strings, default: `[merge]`) - pull request merge options, could be one of:
   - `merge` - enable merge commits PR merging
   - `rebase` - enable rebase PR merge
   - `squash` - enable squash PR merge
 - `deleteBranchOnMerge` (bool) - enables delete branch options on PR merge
 - `topics` (list of strings) - repository topics, keywords in GitHub page description
 - `archived` (bool, readonly) - true if repsitory is archived
 - `disabled` (bool, readonly) - true if repository is disabled
 - `license` (string) - license GitHub key, e.g. (`mit`)
 - visibility (string. default: `"public"`) - one of:
   - `public` - repository is public
   - `private` - repository is private
 - `features` (list of strings) - enables repository features:
   - `issues` - enable issues
   - `wiki` - enable wiki pages
   - `pages` - enable GitHub pages
   - `projects` - enable repository project board
   - `downloads` - ???

Repository metadata must specify repository `name`, and may have `owner`, in case if `owner` is not a current user (current user = token owner). When updating existing repository, metadata must contain `id` of the repository. The full metadata could be fetched with `get` command.

Example:
```yaml
version: v2.0-alpha
kind: Repository
metadata:
    name: gitstrap
    owner: g4s8
    id: 12345
spec:
    description: CLI for managing GitHub repositories
    defaultBranch: master
    mergeStrategy:
        - squash
    deleteBranchOnMerge: true
    topics:
        - cli
        - git
        - github
        - webhooks
    license: mit
    visibility: public
    features:
        - issues
        - pages
        - wiki
        - downloads
```

## Organization

Describes GitHub organization, it has:
 - `name` (string) - the shorthand name of the company.
 - `description` (string) - organization description
 - `company` (string) - organization company affiliation
 - `blog` (string) - URI for organization blog
 - `location` (string) - geo location of organization
 - `email` (string) - public email address
 - `twitter` (string) - twitter account username
 - `verified` (bool, readonly) - true if organization was verified by GitHub

Example:
```yaml
version: v2.0-alpha
kind: Organization
metadata:
    name: artipie
    id: 12345
spec:
    name: Artipie
    description: Binary Artifact Management Toolkit
    blog: https://www.artipie.com
    email: team@artipie.com
    verified: true
```

## WebHook

Describes repository or organization web-hook. It has:
 - `url` (string, required) - webhook URL
 - `contentType` (string, required) - one of:
   - `json` - send JSON payloads
   - `form` - send HTTP form payloads
 - `events` (list of strings, required) - list of [GitHub events](https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads) to trigger web-hook
 - `insecureSsl` (bool, default: false) - if true, disable SSL certificate verification
 - `secret` (string, writeonly) - specify secret payload when creating or updating the hook
 - `active` (bool, default: true) - if false, the hook will be disabled but now removed
 - `selector` (mapping, required) - specifies hook selector. It must have only one of two keys, either:
   - `repository` - the name of repository for this hook
   - `organization` - the name of organization for this hook

Metadata:
 - `owner` (string, optional) - may specify the owner of the repository if hook's selector is repository and repository owner is not a current user
 - `id` (number, required on update) - it must be specified to update existing web-hook, if not specified a new hook will be created. It could be fetched with `get` command.

Example:
```yaml
version: v2.0-alpha
kind: WebHook
metadata:
    owner: g4s8
    id: 12345
spec:
    url: http://example.com/hook
    contentType: json
    events:
        - pull_request
    active: true
    selector:
        repository: gitstrap
```

## Readme

Readme could be described by specification with these fields:
 - `selector` (mapping, required)
   - `repository` (string, required) - the name of repository where this readme will be created
 - `title` (string, optional) - The main title in the readme
 - `abstract` (string, optional) - Short abstract about the repository
 - `topics` (array, optional)
   - `heading` (string, required) - Subs-title for topic
   - `body` (string, required) -  Topic content

Metadata:
 - `owner` (string, optional) - could be specified to create readme in organization or another user repo

Exampple:
```yaml
version: v2
kind: Readme
metadata:
    owner: artipie
spec:
    selector:
      repository: conan-adapter
    title: Conan Artipie adapter
    abstract: >
      Conan is a C/C++ repository, this adapter is an SDK for working
      with Conan data and metadata and a HTTP endpoint for the Conan
      repository.
    topics:
      - heading: How does Conan work
        body: TODO
      - heading: How to use Artipie Conan SDK
        body: TODO
      - heading: How to configure and start Artipie Conan endpoint
        body: TODO
```

## Team

Describes GitHub organization's team, it has:
 - `name` (string, required) - team name
 - `description` (string, optional) - team description
 - `privacy` (string, default: secret) - team privacy
 - `permission` (string, readonly) - team permission. Permission is deprecated when creating or editing a team in an org using the new GitHub permission model.

Metadata:
 - `owner` (string, required) - organization to which the team belongs.
 - `name` (string, required on update) - team slug. It must be specified to update existing team, if not specified, gitstrap will try to update by id.
 - `id` (number, required on update) - it must be specified to update existing team if name is not specified. If name and id are not specified a new team will be created. It could be fetched with `get` command.

Example:
```yaml
version: v2
kind: Team
metadata:
    name: example-team
    owner: artipie
    id: 123456
spec:
    name: Example team
    description: Gitstrap example team
    permission: pull
    privacy: closed
```
