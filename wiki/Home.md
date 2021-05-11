TEST

The `gitstrap` project aims to automate routine operations with GitHub such as managing repositories, organizations, teams, web-hooks and other resources. Each resource is represented as a `yaml` specification document that could be fetched from an existing resource or created from scratch. Each document can be updated and the specification resource configuration can be applied using CLI `gitstrap` tool.

The full specification format can be found here: [/specifications](https://github.com/g4s8/gitstrap/wiki/Specifications)

## Tutorial

 1. [Download and install](#download-and-install)
 2. [Configuration](#configuration)
 3. [CLI overview](#cli-overview)
 4. [Basic examples](#basic-examples)

### Download and install

For Linux system you can use download script to get latest binary:
```bash
curl -L https://raw.githubusercontent.com/g4s8/gitstrap/master/scripts/download.sh | sh
```
This script downloads `gistrap` CLI into `./gitstrap/` path of current directory. The binary could be copied to any of `$PATH` directories.

On MacOS it can be installed using HomeBrew:
```bash
brew tap g4s8/.tap https://github.com/g4s8/.tap
brew install g4s8/.tap/gitstrap
```

On any system (including Windows) the CLI binary can be found at [releases page](https://github.com/g4s8/gitstrap#install).

Alternatively, it can be built from sources (you'll need `git`, `make` and `go` installed on the system):
```bash
git clone --depth=1 https://github.com/g4s8/gitstrap.git
cd gitstrap
make build
sudo make install
```

### Configuration

To use `gitstrap` CLI you need GitHub token. Some features depend on optional permissions:
 - `repo` - required
 - `admin:org` - to manage organizations
 - `admin:org_hook` - to manage organization web-hooks
 - `delete_repo` - to delete repositories

A new token can be created here: https://github.com/settings/tokens/new

When a token is generated, it should be placed at `~/.config/gitstrap/github_token.txt` location.

### CLI overview

The `gitstrap` CLI consists of these primary commands:
 - `get` - get GitHub resource, convert it into `yaml` format and print (possible sub-commands: `repo`, `org`, `hooks`, `teams`, `protection`)
 - `create` - create a new resource from spec `yaml` file, and fail if it already exists
 - `apply` - apply resource specification to existing resource or create it if doesn't exist
 - `delete` - delete resource by `yaml` spec

Global options:
 - `--token` - overrides GitHub token from `~/.config/gitstrap/github_token.txt`

Use `gitstrap --help` or `gitstrap <command> --help` for more details.

### Basic examples

Assuming you have installed and configured `gitstrap`, now you can try to create some resources.

Let's start with simple repository: create a new `yaml` file in current directory: `example-repo.yaml` with content (see the meaning of config fields at [spec reference](https://github.com/g4s8/gitstrap/wiki/Specifications)):
```yaml
version: v2.0-alpha
kind: Repository
metadata:
    name: example
spec:
    description: Example repo created with gitstrap
    license: mit
    visibility: public
    features:
        - issues
        - wiki
```
Create a repository with `gitstrap apply -f example-repo.yaml`. On success, it should print that repository was created.

But it's empty for now, so let's add a README file: create a new specification file `example-readme.yaml`:
```yaml
version: v2
kind: Readme
spec:
    selector:
      repository: example
    title: Example repository
    abstract: >
      This is example repository created with gitstrap
```
And add it to repo using `gitstrap create -f example-readme.yaml` (`gitstrap` doesn't support README updating, so only `create` command could be used).
Now you can check this repository at `/example` location under your account: `https://github.com/<username>/example`.

Let's add a webhook to the repo to call our URL on each `git push` in repository, create a new file `example-hook.yaml` with content:
```yaml
version: v2
kind: WebHook
spec:
    url: https://my-domain.com/hook
    contentType: form
    events:
        - push
    selector:
        repository: example
```
And apply it with `gitstrap apply -f example-hook.yaml`.
