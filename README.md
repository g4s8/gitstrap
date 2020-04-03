# Gitstrap

[![Donate via Zerocracy](https://www.0crat.com/contrib-badge/CF7JL4282.svg)](https://www.0crat.com/contrib/CF7JL4282)
[![DevOps By Rultor.com](http://www.rultor.com/b/g4s8/gitstrap)](http://www.rultor.com/p/g4s8/gitstrap)
[![Managed by Zerocracy](https://www.0crat.com/badge/CF7JL4282.svg)](https://www.0crat.com/p/CF7JL4282)

![Github actions](https://github.com/g4s8/gitstrap/workflows/Go/badge.svg)
[![Hits-of-Code](https://hitsofcode.com/github/g4s8/gitstrap)](https://hitsofcode.com/view/github/g4s8/gitstrap)
[![codebeat badge](https://codebeat.co/badges/89bbb569-fba9-4c68-9b21-e2520b59fbeb)](https://codebeat.co/projects/github-com-g4s8-gitstrap-master)

[![GitHub release](https://img.shields.io/github/release/g4s8/gitstrap.svg?label=version)](https://github.com/g4s8/gitstrap/releases/latest)
[![PDD status](http://www.0pdd.com/svg?name=g4s8/gitstrap)](http://www.0pdd.com/p?name=g4s8/gitstrap)
[![License](https://img.shields.io/github/license/g4s8/gitstrap.svg?style=flat-square)](https://github.com/g4s8/gitstrap/blob/master/LICENSE)

This tool automates routine operations when creating new Github repository.
It can create and configure Github repository from `yaml` configuration file.
Gitstrap helps to: 1) create new repository on Github 2) sync with local directory
3) apply templates, such as README with badges, CI configs, LICENSE stuff, etc
4) configure webhooks for Github repo 5) invite collaborators

Also, it has powerfull extensions for integration with external services or with
Github.

## How to use

First you need to install it.

To get binary for your platform use [download script](https://github.com/g4s8/gitstrap/blob/master/scripts/download.sh):
```sh
curl -L https://raw.githubusercontent.com/g4s8/gitstrap/master/scripts/download.sh | sh
```

For Gentoo Linux you can merge it from my repo [Layman](https://wiki.gentoo.org/wiki/Layman) overlay:
```sh
sudo layman -o https://raw.githubusercontent.com/g4s8-overlay/layman/master/repositories.xml -a g4s8
sudo emerge -av dev-vcs/gitstrap
```

On MacOS you can install it using `brew` tool:
```sh
brew tap g4s8/.tap https://github.com/g4s8/.tap
brew install g4s8/.tap/gitstrap
```

Alternatively, you can build it using `go get github.com/g4s8/gitstrap`

### Get GitHub token

Go to settings (profile settings, developer settings, personal access token, generate new token):
https://github.com/settings/tokens/new
and select all `repo` checkboxes and `delete_repo` checkbox (in case you want gitstrap to be able to
delete repositories). You may use this token as CLI option (`gitstrap --token=ABCD123 apply`)
or save it in `~/.config/gitstrap/github_token.txt` file.

### Create new project

Create new config YAML file with gitstrap configuration:
```yaml
# test.yaml
---
version: v2
github:
  # repository name
  name: test1234
  # repository description
  description: just a test repo
  # private repository (optional, public by default)
  private: true
  # webhooks
  hooks:
    - url: "https://myserver.com/hooks"
      type: form
      events:
        - push
      active: true
  # add collaborators
  collaborators:
    - rultor
extensions:
  # generate README file and upload it
  readme:
    # Readme title
    title: "Test project"
    # List of badges
    badges:
      - alt: CircleCI
        img: https://circleci.com/gh/g4s8/gitstrap.svg?style=svg
        link: https://circleci.com/gh/g4s8/gitstrap
  # Enable https://www.0pdd.com/
  0pdd:
    verbose: false
```

Apply this config: `gitstrap apply -f test.yaml`

### Update the project

Later, you may decide to update this configuration, e.g. add new webhook,
removed collaborators or change the description. Just edit your YAML config file
and run `gitstrap apply -f ./filename.yaml` again.

### Delete repo

To delete GitHub repository use `delete` command:
`gitstrap delete -f config.yaml`.

### Obtain config from existing repo

It's possible to generate yaml config from existing repo.
It may be usefull if you want to use this repository as a template for new one,
or you want to change some details via `gitstrap` tool.<br/>
**Important: extensions may not recognize repository configuration, only primary config sections
will be imported.**<br/>
Run: `gitstrap get -o yaml reponame`, or redirect it to file if you want to edit it:
`gitstrap get -o yaml reponame > config.yaml`.


### Manage organization repos

To manage organization repositories, use `--org=orgname` option, e.g. to create new repo
from `example.yaml` config in organization `test`, use `gitstrap --org=test apply -f example.yaml`.
Make sure you have permissions to create or delete repositories in this organization.

## Configuration

Current supported version is `v2`. Gitstrap can support `v1` as well, but
`v1` doesn't include all features. Include this line in configuration root:
```yaml
version: v2
```

There are two main sections of gitstrap configuration:
 - `github` - primary GitHub configuration: name, hooks, users, etc
 - `extensions` - additional extensions to intergrate with external services
 or simplify common complex routines, e.g. readme extension

### GitHub configuration

GitHub primary configuration options:
 - `name` - (required) repository name
 - `description` - (optional) repository description
 - `private` - (optional, default false) public repository or private
 - `hooks` - (optional) list of webhooks
 - `collaborators` - (optional) list of usernames of collaborators

Each webhooks has these options:
 - `url` - (required) hook URL
 - `type` - (required) string: `form` or `json`
 - `events` - (required) list of events strings (or `*` for all), see the list of
 [supported events](https://developer.github.com/webhooks/#events)
 - `active` - (optional, true by default) use `false` to disable the webhook but not
 delete it

`url` paramter should be unique across all webhooks.

Webhooks example:
```yaml
hooks:
  - url: "https://myserver.com/hooks"
    type: json
    events:
      - '*'
    active: true
  - url: http://www.0pdd.com/hook/github
    type: form
    events:
      - push
    active: false
  - url: https://notify.travis-ci.org
    type: form
    events:
      - create
      - delete
      - issue_comment
```

Collaborators example:
```yaml
collaborators:
  - johndoe
  - janedoe
```

### Extensions
Currently supported extensions: `readme`, `0pdd`.

Readme extension generates new `README.md` file if doesn't exist in the repo
by specified paramters:
 - `title` - readme title, optional
 - `header` - heading text, optional
 - `badges` - list of badges, where each badge must include parameters:
   - `alt` - alternative text to display
   - `img` - image URL to display
   - `link` - link to open on click

Example:
```yaml
extensions:
  readme:
    title: "Test project"
    badges:
      - alt: CircleCI
        img: https://circleci.com/gh/g4s8/gitstrap.svg?style=svg
        link: https://circleci.com/gh/g4s8/gitstrap
```

0pdd extension configures https://www.0pdd.com service, it creates new
`.pdd` configuration file and upload to the repo if doesn't exist,
it appends 0pdd webhook and invites 0pdd collaborator to create new tickets:
```yaml
extensions:
  0pdd:
    verbose: true
    minWords: 20
```

