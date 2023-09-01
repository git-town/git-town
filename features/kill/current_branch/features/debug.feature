Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "current"
    And a feature branch "other"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |

  # TODO: remove redundant "git config -lz --local"
  Scenario: result
    When I run "git-town kill --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git version                                       |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git remote                                        |
      |         | backend  | git status                                        |
      |         | backend  | git rev-parse --abbrev-ref HEAD                   |
      | current | frontend | git fetch --prune --tags                          |
      |         | backend  | git branch -vva                                   |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git status --porcelain --ignore-submodules        |
      | current | frontend | git push origin :current                          |
      |         | frontend | git checkout main                                 |
      |         | backend  | git rev-parse current                             |
      |         | backend  | git log main..current                             |
      | main    | frontend | git branch -D current                             |
      |         | backend  | git config --unset git-town-branch.current.parent |
      |         | backend  | git show-ref --quiet refs/heads/other             |
      |         | backend  | git show-ref --quiet refs/heads/current           |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git checkout other                                |
      |         | backend  | git checkout main                                 |
    And it prints:
      """
      Ran 22 shell commands.
      """
    And the current branch is now "main"
