To work on the project you need to have `git`, `go` and `make` tools.

To start working:
 1. Fork the repo
 2. Clone your fork
 3. Create a branch for changes
 4. Commit changes to your local branch
 5. Push to fork
 6. Create pull request to upstream

To build the project use `make` command from project root:
 - `make build` - builds the binary
 - `make test` - starts unit tests
 - `make link` - starts linter

Before commiting changes make sure that `make` command passes successfully.

Project layout is:
 - `internal/spec/` - specification definitions, see more details in [wiki](https://github.com/g4s8/gitstrap/wiki/Specifications)
 - `internal/gitstrap` - all supported operations on specificaitons
 - `cmd/gitstrap` - CLI main

### Debugging

To debug it on your own repository you need to [generate a token](https://github.com/g4s8/gitstrap#get-github-token).

For debugging, use `DEBUG=1` environment variable, e.g. `DEBUG=1 ./gitstrap get repo gitstrap`

## Code style

Code style is checked by [golangci-lint](https://golangci-lint.run/) tool, you may need to install it if don't have. Run `mkae lint` to ensure
your code style is OK. In pull requests it checked automatically, and CI workflow block PR from merging.

## Pull request style

Primary PR rule: it's the responsibility of PR author to bring the changes to the master branch.

Other important mandatory rule - it should refer to some ticket. The only exception is a minor type fix in documentation.

Pull request should consist of two mandatory parts:
 - "Title" - says **what** is the change, it should be one small and full enough sentence with only necessary information
 - "Description" - says **how** this pull request fixes a problem or implements a new feature

### Title

Title should be as small as possible but provide full enough information to understand what was done (not a process),
and where from this sentence.
It could be capitalized If needed, started from capital letter, and should not include links or references
(including tickets numbers).

Good PR titles examples:
 - Fixed Maven artifact upload - describes what was done: fixed, the what was the fixed: artifact upload, and where: Maven
 - Implemented GET blobs API for Docker - done: implemented, what: GET blobs API, where: Docker
 - Added integration test for Maven deploy - done: added, what: integration test for deploy, where: Maven

Bad PR titles:
 - Fixed NPE - not clear WHAT was the problem, and where; good title could be: "Fixed NPE on Maven artifact download"
 - Added more tests - too vague; good: "Added unit tests for Foo and Bar classes"
 - Implementing Docker registry - the process, not the result; good: "Implemented cache layer for Docker proxy"

### Description

Description starts with a ticket number prefixed with one of these keywords: (Fixed, Closes, For, Part of),
then a hyphen, and a description of the changes.
Changes description provides information about **how** the problem from title was fixed.
It should be a short summary of all changes to increase readability of changes before looking to code,
and provide some context. The format is
`(<keyword>For|Closes|Fixes|Part of) #(<ticket>\d+) - (<details>.+)`,
e.g.: `For #123 - check if the file exists before accessing it and return 404 code if doesn't`.

Good description describes the solution provided and may have technical details, it isn't just a copy of the title.
Examples of good descriptions:
 - Added a new class as storage implementation over S3 blob-storage, implemented `value()` method, throw exceptions on other methods, created unit test for value
 - Fixed FileNotFoundException on reading blob content by checking if file exists before reading it. Return 404 code if doesn't exist

### Merging

We merge PR only if all required CI checks passed and after approval of repository maintainers.
We merge using squash merge, where commit messages consists of two parts:
```
<PR title>

<PR description>
PR: <PR number>
```
GitHub automatically inserts title and description as commit messages, the only manual work is a PR number.

### Review

It's recommended to request review from `@artipie/contributors` if possible.
When the reviewers starts the review it should assign the PR to themselves,
when the review is completed and some changes are requested, then it should be assigned back to the author.
On approve: if reviewer and repository maintainer are two different persons,
then the PR should be assigned to maintainer, and maintainer can merge it or ask for more comments. 

The workflow:
```
<required> (optional)
        PR created |   Review   | Request changes | Fixed changes | Approves changes | Merge |
assignee: <none>  -> <reviewer> ->    (author)    ->  (reviewer)  ->   <maintainer>  -> <none>
```

When addressing review changes, two possible strategies could be used:
 - `git commit --ammend` + `git push --force` - in case of changes are minor or obvious, both sides agree
 - new commit - in case if author wants to describe review changes and keep it for history,
 e.g. if author doesn't agree with reviewer or maintainer, he|she may want to point that this changes was
 asked by a reviewer. This commit is not going to the master branch, but it will be linked into PR history.

### Commit style

Commit styles are similar to PR, PR could be created from commit message: first line goes to the title,
other lines to description:
```
Commit title - same as PR title

For #123 - description of the commit goes
to PR description. It could be multiline `and` include
*markdown* formatting.
```
