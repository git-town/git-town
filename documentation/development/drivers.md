# Drivers

_The following refers to the commands `git-new-pull-request` and `git-repo`.

_Drivers_ implement third-party specific functionality in a standardized way.
For example, the [GitHub driver](./src/drivers/code_hosting/github.sh)
implements GitHub-related operations like creating a pull request there.

There is also an analogous
[Bitbucket driver](./src/drivers/code_hosting/bitbucket.sh)
that does the same things on Bitbucket.
Both drivers are part of the [code hosting](./src/drivers/code_hosting) _driver family_.

The functions that a driver needs to implement are described in the
documentation for the respective driver family.

In order to use a driver, a script simply needs to activate the respective
driver family.
The driver family's activation script then automatically determines
the appropriate driver for the current environment and runs it.
