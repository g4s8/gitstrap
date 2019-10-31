To contribute you have to join [Zerocracy](https://www.0crat.com)
and get `DEV` or `REV` role in this project: [`CF7JL4282`](https://www.0crat.com/p/CF7JL4282).

You need to read and understand the
[Policy](https://www.zerocracy.com/policy.html) before contributing.

Git guidelines:
 - start all commit message with tiket number, e.g.: `#1 - some message`
 (to avoid issues with comment char, check this answer: https://stackoverflow.com/a/14931661/1723695)
 - try to describe shortly your changes in commit message, long description
 should be provided in PR body
 - name your branches starting with ticket number, e.g. `1-some-bug` branch for #1 ticket
 - use `merge`, not `rebase` when merging changes from `master` to local branch
 - one commit per change, don't rebase all commits into single one for PR
 - avoid `push --force` where possible, it can be used in rare cases, e.g. if
 you pushed binary file by mistake, then you can remove it with `push --force`
 or if you mistyped ticket number or message, you can fix it with `push --force`

To change the code folow these steps:
 1. Fork & clone the repo
 2. Make changes
 3. Make sure code can be built and all tests passed: `go build . && go test .`
 4. Make sure all `gometalinter` checks passed: `gometalinter .`
 5. Submit a pull request, make sure all PR checks passed (green circles)
 6. Wait for review and apply review changes if asked
 7. Done - your code now in master and will be available with next release
