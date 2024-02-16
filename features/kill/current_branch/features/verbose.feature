Feature: display all executed Git commands

  Background:
    Given the current branch is a feature branch "current"
    And a feature branch "other"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And the current branch is "current" and the previous branch is "other"

  Scenario: result
    When I run "git-town kill --verbose"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git version                                       |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git stash list                                    |
      |         | backend  | git remote                                        |
      |         | backend  | git status --long --ignore-submodules             |
      |         | backend  | git rev-parse --abbrev-ref HEAD                   |
      | current | frontend | git fetch --prune --tags                          |
      |         | backend  | git branch -vva                                   |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git status --long --ignore-submodules             |
      | current | frontend | git push origin :current                          |
      |         | frontend | git checkout other                                |
      |         | backend  | git log main..current                             |
      | other   | frontend | git branch -D current                             |
      |         | backend  | git config --unset git-town-branch.current.parent |
      |         | backend  | git show-ref --verify --quiet refs/heads/current  |
      |         | backend  | git branch -vva                                   |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git stash list                                    |
    And it prints:
      """
      Ran 22 shell commands.
      """
    And the current branch is now "other"
