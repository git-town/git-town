# Git Town Release Notes

## 0.4
* pushes tags to the remote when running "git sync" on the main branch
* added support for non-feature branches (cannot be shipped and do not merge main when syncing)
* git ship merges main into the feature branch before squash merging
* updated output to show each git command and its output, updated error messages
* skips unnecessary pushes


## 0.3
* multi-user support for feature branches (https://github.com/Originate/git-town/issues/35)
* git sync-fork
* stores configuration in the Git configuration instead of a dedicated file
* only makes one fetch from the central repo per session
* automatically prunes remote branches when fetching updates
* always cleans up abort and continue scripts after using one of them
* simpler readme, dedicated RDD document
* <a href="http://cukes.info" target="_blank">Cucumber</a> feature specs (you need Ruby 2.x)
* much faster testing thanks to completely local test Git repos (https://github.com/Originate/git-town/issues/25)


## 0.2.2
* fixes "unary" error messages
* lots of output and documentation improvements


## 0.2.1
* better terminal output
* Travis CI improvements
* better documentation


## 0.2
* displays the duration of specs
* only pulls the main branch if it has a remote
* --abort options to abort failed Git Town operations
* --continue options to continue some Git Town operations after fixing the underlying issues
* can be installed through Homebrew
* colored test output
* display summary after tests
* exit with proper status codes
* better documentation


## 0.1
* git hack, git sync, git extract, git ship
* basic test framework
* Travis CI integration
* self-hosting: uses Git Town for Git Town development

