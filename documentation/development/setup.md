## Developer machine setup

You need to have these things running on your computer. Please refer to their
websites for help getting them installed.

- [Go](https://golang.org) version 1.9 or higher
- [Yarn](https://yarnpkg.com/)
- [Make](https://www.gnu.org/software/make)
  - Mac and Linux users should be okay out of the box
  - Windows users should install
    [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
- [scc](https://github.com/boyter/scc)

Fork Git Town and clone your fork into a directory outside your GOPATH. Git Town
uses Go modules and doesn't work properly inside the GOPATH.

Cd into the directory you just cloned and run
<code textrun="verify-make-command">make setup</code> to download additional
tooling and dependencies.

To make sure everything works,

- build the tool: <code textrun="verify-make-command">make build</code>
  - now you have `$GOPATH/bin/git-town` compiled from your local source code
- run the tests: <code textrun="verify-make-command">make test</code>
