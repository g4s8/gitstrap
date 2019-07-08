# Gitstrap

[![Donate via Zerocracy](https://www.0crat.com/contrib-badge/CF7JL4282.svg)](https://www.0crat.com/contrib/CF7JL4282)
[![DevOps By Rultor.com](http://www.rultor.com/b/g4s8/gitstrap)](http://www.rultor.com/p/g4s8/gitstrap)
[![Managed by Zerocracy](https://www.0crat.com/badge/CF7JL4282.svg)](https://www.0crat.com/p/CF7JL4282)

[![GitHub release](https://img.shields.io/github/release/g4s8/gitstrap.svg?label=version)](https://github.com/g4s8/gitstrap/releases/latest)
[![Build Status](https://img.shields.io/travis/g4s8/gitstrap.svg?style=flat-square)](https://travis-ci.org/g4s8/gitstrap)
[![CircleCI](https://circleci.com/gh/g4s8/gitstrap.svg?style=svg)](https://circleci.com/gh/g4s8/gitstrap)
[![Hits-of-Code](https://hitsofcode.com/github/g4s8/gitstrap)](https://hitsofcode.com/view/github/g4s8/gitstrap)

[![PDD status](http://www.0pdd.com/svg?name=g4s8/gitstrap)](http://www.0pdd.com/p?name=g4s8/gitstrap)
[![License](https://img.shields.io/github/license/g4s8/gitstrap.svg?style=flat-square)](https://github.com/g4s8/gitstrap/blob/master/LICENSE)

This tool automates routine operations when creating new Github repository.
It can create and configure Github repository from `yaml` configuration file.
Gitstrap helps to: 1) create new repository on Github 2) sync with local directory
3) apply templates, such as README with badges, CI configs, LICENSE stuff, etc
4) configure webhooks for Github repo 5) invite collaborators

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

*Before using, make sure you have running `ssh-agent` daemon with imported ssh key for Github
to be able to push, fetch and pull.*

To bootstrap new repository:
 1. Create new directory for your project
 2. Write `.gitstrap.yaml` config in project root
 3. Create Github API token with repo permissions (if doesn't exist, need's to be done once)
 4. Run `gitstrap -token=your-api-token create`

To remove repository (keep source files) run `gitstrap -token=your-api-token destroy`

To create or destroy a repo for an organization make sure you have permissions to create
or delete repositories and use `-org=orgname` gitstrap option.

## Configuration
The default configuration file is `.gitstrap.yaml`, you can specify another file with `-config=my-config.yaml` option.

There's sample config yaml:
```yaml
gitstrap:
    # gitstrap config version, should be v1
    version: v1
    github:
        repo:
            # github repository name (optional, current directory name if empty)
            name: gitstrap
            # github repository description
            description: "Command line tool to bootstrap Github repository"
            # true if private (optional, default false) 
            private: false
            # github webhooks
            hooks:
                  # webhook url
                - url: "http://p.rehttp.net/http://www.0pdd.com/hook/github"
                  # webhook type: form or json
                  type: form
                  # events to send (see github docs)
                  events:
                      - push
                  # false to create webhook in inactive state (optional, default true) 
                  active: true
            # github logins to add as collaborators
            collaborators:
                - "rultor"
                - "0pdd"
    # (optional) templates to apply, see Templates README section
    templates:
          # file name in repository
        - name: "README.md"
          # template url
          url: "https://raw.githubusercontent.com/g4s8/gitstrap/master/templates/README.md"
        - name: "LICENSE"
          location: "/home/g4s8/.gitstrap/LICENSE.mit"
    # (optional) these params can be accessed from template, just a key-value pairs
    params:
        rultor: true
        travis: true
        readmeContrib: |
            Fork repository, clone it, make changes,
            push to new branch and submit a pull request.
        pdd: true
        license: MIT
```

## Templates
You can create any file templates to apply them using project configuration.
Templates can use [golang template](https://golang.org/pkg/text/template/) patterns.

*Example: README.md*
```markdown
# {{.Repo.Name}}.
{{.Repo.Description}}

{{if .Gitstrap.Params.eoPrinciples }}[![EO principles respected here](http://www.elegantobjects.org/badge.svg)](http://www.elegantobjects.org){{end}}
{{if .Gitstrap.Params.rultor}}[![DevOps By Rultor.com](http://www.rultor.com/b/{{.Repo.Owner.Login}}/{{.Repo.Name}})](http://www.rultor.com/p/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}

{{if .Gitstrap.Params.travis}}[![Build Status](https://img.shields.io/travis/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://travis-ci.org/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.appveyor}}[![Build status](https://ci.appveyor.com/api/projects/status/{{.Gitstrap.Params.appveyor}}?svg=true)](https://ci.appveyor.com/project/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.pdd}}[![PDD status](http://www.0pdd.com/svg?name={{.Repo.Owner.Login}}/{{.Repo.Name}})](http://www.0pdd.com/p?name={{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.license}}[![License](https://img.shields.io/github/license/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://github.com/{{.Repo.Owner.Login}}/{{.Repo.Name}}/blob/master/LICENSE){{end}}
{{if .Gitstrap.Params.codecov}}[![Test Coverage](https://img.shields.io/codecov/c/github/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://codecov.io/github/{{.Repo.Owner.Login}}/{{.Repo.Name}}?branch=master){{end}}
```
this template uses [`.gitstrap.yaml`](https://github.com/g4s8/gitstrap/blob/master/.gitstrap.yaml) config and produces `README.md`:
```markdown
# gitstrap.
Command line tool to bootstrap Github repository

[![DevOps By Rultor.com](http://www.rultor.com/b/g4s8/gitstrap)](http://www.rultor.com/p/g4s8/gitstrap)

[![Build Status](https://img.shields.io/travis/g4s8/gitstrap.svg?style=flat-square)](https://travis-ci.org/g4s8/gitstrap)

[![PDD status](http://www.0pdd.com/svg?name=g4s8/gitstrap)](http://www.0pdd.com/p?name=g4s8/gitstrap)
[![License](https://img.shields.io/github/license/g4s8/gitstrap.svg?style=flat-square)](https://github.com/g4s8/gitstrap/blob/master/LICENSE)
```
