# Drivers

_The following refers to the commands `git-new-pull-request` and `git-repo`._

_Drivers_ implement third-party specific functionality in a standardized way.
For example, the [GitHub driver](/src/drivers/github.go)
implements GitHub-related operations like creating a pull request there.

There is also an analogous
[Bitbucket driver](/src/drivers/bitbucket.go)
that does the same things on Bitbucket.
Both drivers are part of the [code hosting](/src/drivers/code_hosting_driver.go) _driver family_.

The functions that a driver needs to implement are described in the
documentation for the respective driver family.

In order to use a driver, a script simply needs to activate the respective
driver family.
The driver family's activation script then automatically determines
the appropriate driver for the current environment and runs it.
