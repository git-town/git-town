## Developer machine setup

You need to have these things running on your computer. Please refer to their
websites for help getting them installed.

- [Go](https://golang.org) version 1.9 or higher
- [Ruby 2.4.1](https://www.ruby-lang.org/en/documentation/installation) to run
  the legacy feature tests
- [Yarn](https://yarnpkg.com/)
- [Make](https://www.gnu.org/software/make)
  - Mac and Linux users should be okay out of the box
  - Windows users should install
    [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

Now clone Git Town into `$GOPATH/src/github.com/git-town/git-town`. Until we
have enabled Go modules, Git Town must live inside the GOPATH.

Then, cd into the directory you just cloned, and run
<code textrun="verify-make-command">make setup</code> to download additional
tooling and dependencies.

To make sure everything works:

- build the tool: <code textrun="verify-make-command">make build</code>
  - now you have `$GOPATH/bin/git-town` compiled from your local source code
