# Gitstrap

![Github actions](https://github.com/g4s8/gitstrap/workflows/Go/badge.svg)
[![Hits-of-Code](https://hitsofcode.com/github/g4s8/gitstrap)](https://hitsofcode.com/view/github/g4s8/gitstrap)
[![codebeat badge](https://codebeat.co/badges/89bbb569-fba9-4c68-9b21-e2520b59fbeb)](https://codebeat.co/projects/github-com-g4s8-gitstrap-master)

[![GitHub release](https://img.shields.io/github/release/g4s8/gitstrap.svg?label=version)](https://github.com/g4s8/gitstrap/releases/latest)
[![PDD status](http://www.0pdd.com/svg?name=g4s8/gitstrap)](http://www.0pdd.com/p?name=g4s8/gitstrap)
[![License](https://img.shields.io/github/license/g4s8/gitstrap.svg?style=flat-square)](https://github.com/g4s8/gitstrap/blob/master/LICENSE)

Manage your GitHub repositories as a set of resouce configuration files!

Gitstrap automates routine operations with Github.
It can create and configure Github repositories, teams, readmes, organizations, etc
from `yaml` specification files.
It helps to: 1) create new repository on Github 2) manage repositories permissions
3) keep all organization repositories configuration in yaml files in one directory
4) configure webhooks for Github repo 5) configure branch protection rules
6) many others


See Wiki for [full documentation](https://github.com/g4s8/gitstrap/wiki/Specifications).

# Resource specifications

Gitstrap works with specification documents. Each document is a YAML mapping with a format:
```yaml
# specification version (now only v2.0-alpha is supported)
version: v2.0-alpha
# kind of specification
kind: Repository
# metadata, could be different depends on kind
metadata:
    # resource name
    name: gitstrap
    # resource owner
    owner: g4s8
    # resource id
    id: 164098671
    # optional annotations
    annotations: {}
# specification declaration
spec: {}
```

At the moment, Gitstrap supports these specifications:
 - `Repository` - GitHub repository specs
 - `Readme` - README file in repository

## Repo-spec

Repository specification allows to manage GitHub repository with `gitstrap` CLI.
It consists of these parameters:
```yaml
spec:
    # repository description
    description: CLI for managing GitHub repositories
    # default branch name
    defaultBranch: master
    # merge strategies (commit, rebase, squash)
    mergeStrategy:
        - squash
    # enables delete branch on merge options
    deleteBranchOnMerge: true
    # repository topics (keywords)
    topics:
        - cli
        - git
        - github
        - webhooks
    # license key
    license: mit
    # repository visibility (private or public)
    visibility: public
    # features enabled (issues, wiki, pages, projects, downloads)
    features:
        - issues
        - pages
        - downloads
    # repository collaborators permissions (available only for repo owner)
    collaborators:
        # collaborator name
        - name: g4s8
          # permissions (admin, push, pull)
          permissions:
            - admin
            - push
            - pull
```

Repository spec supports these metadata params:
 - `name` - repository name
 - `owner` - repository owner (empty for token owner)
 - `id` - repository ID for updating existing

# Quickstart

 1. Download `gitstrap` CLI (see [Install](#install) section)
 2. Get configuration from any of your repositories or from this one: `gitstrap get --owner=g4s8 gitstrap > repo.yaml`
 3. Edit YAML config (see [Specification](https://github.com/g4s8/gitstrap/wiki/Specifications) reference)
 4. Create or update you repository with `gitstrap apply -f repo.yaml`


## Install

First you need to install it.

To get binary for your platform use [download script](https://github.com/g4s8/gitstrap/blob/master/scripts/download.sh):
```sh
curl -L https://raw.githubusercontent.com/g4s8/gitstrap/master/scripts/download.sh | sh
```

On MacOS you can install it using `brew` tool:
```sh
brew tap g4s8/.tap https://github.com/g4s8/.tap
brew install g4s8/.tap/gitstrap
```

Alternatively, you can build it using `go get github.com/g4s8/gitstrap`

## Get GitHub token

To use `gitstrap` you need GitHub token.
Go to settings (profile settings, developer settings, personal access token, generate new token):
https://github.com/settings/tokens/new
and select all `repo` checkboxes and `delete_repo` checkbox (in case you want gitstrap to be able to
delete repositories). You may use this token as CLI option (`gitstrap --token=ABCD123 apply`)
or save it in `~/.config/gitstrap/github_token.txt` file.

