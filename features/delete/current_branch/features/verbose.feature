Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the current branch is "current"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And the current branch is "current" and the previous branch is "other"

  Scenario: result
    When I run "git-town delete --verbose"
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git version                                       |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git config -lz --includes --global                |
      |         | backend  | git config -lz --includes --local                 |
      |         | backend  | git status --long --ignore-submodules             |
      |         | backend  | git remote                                        |
      |         | backend  | git rev-parse --abbrev-ref HEAD                   |
      | current | frontend | git fetch --prune --tags                          |
      |         | backend  | git stash list                                    |
      |         | backend  | git branch -vva --sort=refname                    |
      |         | backend  | git remote get-url origin                         |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git remote get-url origin                         |
      | current | frontend | git push origin :current                          |
      |         | frontend | git checkout other                                |
      | other   | frontend | git branch -D current                             |
      |         | backend  | git config --unset git-town-branch.current.parent |
      |         | backend  | git show-ref --verify --quiet refs/heads/current  |
      |         | backend  | git branch -vva --sort=refname                    |
      |         | backend  | git config -lz --includes --global                |
      |         | backend  | git config -lz --includes --local                 |
      |         | backend  | git stash list                                    |
    And Git Town prints:
      """
      Ran 22 shell commands.
      """
    And the current branch is now "other"
