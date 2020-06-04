# Code Hosting Drivers

Code hosting drivers allow commands like `git-new-pull-request`, `git-repo`, and
`git ship` to communicate with the API of your hosting service.

Drivers implement the [CodeHostingDriver](/src/drivers/core.go) interface.
Driver implementations are available for [GitHub](/src/drivers/github.go),
[Gitea](/src/drivers/gitea.go), [Bitbucket](/src/drivers/bitbucket.go), and
[GitLab](/src/drivers/bitbucket.go). To use a driver, call `drivers.Load()`.
