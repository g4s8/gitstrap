# {{.Repo.Name}}.
{{.Repo.Description}}

{{if .Gitstrap.Params.eoPrinciples }}[![EO principles respected here](http://www.elegantobjects.org/badge.svg)](http://www.elegantobjects.org){{end}}
{{if .Gitstrap.Params.rultor}}[![DevOps By Rultor.com](http://www.rultor.com/b/{{.Repo.Owner.Login}}/{{.Repo.Name}})](http://www.rultor.com/p/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}

{{if .Gitstrap.Params.travis}}[![Build Status](https://img.shields.io/travis/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://travis-ci.org/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.appveyor}}[![Build status](https://ci.appveyor.com/api/projects/status/{{.Gitstrap.Params.appveyor}}?svg=true)](https://ci.appveyor.com/project/{{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.pdd}}[![PDD status](http://www.0pdd.com/svg?name={{.Repo.Owner.Login}}/{{.Repo.Name}})](http://www.0pdd.com/p?name={{.Repo.Owner.Login}}/{{.Repo.Name}}){{end}}
{{if .Gitstrap.Params.license}}[![License](https://img.shields.io/github/license/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://github.com/{{.Repo.Owner.Login}}/{{.Repo.Name}}/blob/master/LICENSE){{end}}
{{if .Gitstrap.Params.codecov}}[![Test Coverage](https://img.shields.io/codecov/c/github/{{.Repo.Owner.Login}}/{{.Repo.Name}}.svg?style=flat-square)](https://codecov.io/github/{{.Repo.Owner.Login}}/{{.Repo.Name}}?branch=master){{end}}

{{if .Gitstrap.Params.readmeContrib}}
## Contribution
{{.Gitstrap.Params.readmeContrib}}
{{end}}

