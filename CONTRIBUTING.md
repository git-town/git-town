# How to contribute

Git Town is a welcoming community, and we'd love for everyone to bring
their contributions to make it even better.
We appreciate contributions of any size.

* Found a bug or have an idea for a new feature? - [Open an issue](https://github.com/Originate/git-town/issues/new)
* Fixed a bug or created a new feature that others will enjoy? - [Create a pull request](https://help.github.com/articles/using-pull-requests/)

This guide will help you get started and outline some things you should know when developing Git Town.


## Requirements

* Ruby 2.2
  (install [directly](https://www.ruby-lang.org/en/documentation/installation),
  or via a ruby manager like [rvm](https://rvm.io/)
  or [rbenv](https://github.com/sstephenson/rbenv))
  * language the tests are written in
* [ShellCheck](https://github.com/koalaman/shellcheck)
  * used in the linting process to find common errors in the Bash code


## Setup

* install the [requirements](#requirements)
* fork and clone the repository to your machine
* run `bundle` to install ruby gems
* optionally run `rake` to make sure all tests pass on your machine


## Testing

* tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* all features need to have comprehensive test coverage
* source code and test files must pass the linters

```bash
# running the different test types
rake         # runs all tests
rake lint    # runs the linters
rake test    # runs the feature tests

# running individual scenarios/features
cucumber <filename>[:<lineno>]
cucumber -n '<scenario/feature name>'

# running several features in parallel
bin/cuke [cucumber parameters]

# auto-fixing formatting issues
rake format  # Run formatters (fixes some lint errors)
```

The `rake [parameters]` commands above can also be run as `bundle exec rake [parameters]`
if you encounter issues.

Git Town's [CI server](https://circleci.com/gh/Originate/git-town)
automatically tests all commits and pull requests,
and notifies you via email and through status badges in pull requests
about problems.


## Debugging

To see the output of the Git commands run in tests, you can set the
`DEBUG_COMMANDS` environment variable while running your specs:

```bash
$ DEBUG_COMMANDS=true cucumber <filename>[:<lineno>]
```

Alternatively, you can also add a `@debug-commands` flag to the respective
Cucumber spec:

  ```cucumber
  @debug-commands
  Scenario: foo bar baz
    Given ...
  ```

For even more detailed output, you can use the `DEBUG` variable or tag
in a similar fashion.
If set, Git Town prints every shell command executed during the tests
(includes setup, inspection of the Git status, and the Git commands),
and the respective console output.


## Pull Requests

Each pull request (PR) should have the same (optional) description that it will
have when committed later and include the
[issue](https://github.com/Originate/git-town/issues) it resolves.

When merging approved PRs:
* use `git ship`
* the message for the squashed commit should follow the
  [formatting guidelines for commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
* [mention the closed issue](https://help.github.com/articles/closing-issues-via-commit-messages)
in the commit body, so that the respective issue is automatically closed.

Example of a commit message:

```
Automatically foo commits in "git bar"

Fooing changes before they get barred into a separate branch
keeps the final foo simpler when shipping that branch later.

Implements #123
```


## Architecture

_The following refers to all commands except `git-town`._

Each Git Town command begins by inspecting the current state of the Git repository
(which branch you are on, whether you have open changes).
If there are no errors, it generates a list of steps to run.
Each step is a bash function that wraps an individual Git command.
This list is then executed one by one.

For discussion around this architecture see
[#199](https://github.com/Originate/git-town/issues/199),
where it was proposed.


### Drivers

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


### Branch Hierarchy Architecture

Since code reviews can take a while,
many developers work on several features in parallel.
These features often depend on each other.
To support this common use case, Git Town provides an hierarchical branching model
that is more opinionated than the very generic branching of vanilla Git.
In Git Town's world, feature branches can be "children" of other feature branches.

As an example, lets assume a repo with the following setup:

```
-o--o-- master
  \
   o--o--o-- feature1
       \
        o-- feature2
```

In this example, feature 1 (which was cut straight from the master branch) is currently under review.
While waiting for the LGTM there, the developer has started to work on the next feature.
This work (let's call it "feature 2") needs some of the changes that are introduced by feature 1.
Since feature 1 hasn't shipped yet, we can't cut feature 2 straight off master,
but must cut it off feature 1, so that feature 2 sees the changes made by feature 1.

This means the feature branch `feature1` is cut directly from `master`,
and `feature2` is cut from `feature1`, making it a child branch of `feature1`.

This "ancestry line" of branches is preserved at all times,
and impacts a lot of Git Town's commands.
For example, child branches cannot be shipped before their parents.
When syncing, Git Town syncs the parent branch first,
then merges the parent branch into its children branches.
When creating a pull request for `feature2`,
Git Town only displays the changes between `feature2` and `feature1`,
not the diff against `master`.

Git Town stores the information about this branch hierarchy in the Git configuration for the repo.
Two types of keys are used for this. The first one is __git-town.branches.parent__.
It lists which branch is the immediate parent branch of the given branch.
```
git-town.branches.parent.feature1=master
git-town.branches.parent.feature2=feature1
```

Git Town also caches the full ancestral line of each feature branch, top-down,
in a key called __git-town.branches.parents__:
* `git-town.branches.parents.feature2=master,feature1`
  lists that in order to sync `feature2`, we need to first update `master`,
  then merge master into `feature1`, then `feature1` into `feature2`.


## Documentation

Every Git Town command
* has a [man page](./man/man1)
* has a [Markdown page](./documentation/commands) that is identical to the man page
* is listed on the [git-town man page](./man/man1/git-town.1)
* is listed on the [README](./README.md)
