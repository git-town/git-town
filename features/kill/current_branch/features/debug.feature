Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "current"
    And a feature branch "other"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |

  Scenario: result
    When I run "git-town kill --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git rev-parse                                     |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git version                                       |
      |         | backend  | git branch -a                                     |
      |         | backend  | git status                                        |
      |         | backend  | git rev-parse --abbrev-ref HEAD                   |
      |         | backend  | git branch                                        |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git remote                                        |
      | current | frontend | git fetch --prune --tags                          |
      |         | backend  | git branch -r                                     |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git status --porcelain --ignore-submodules        |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git status --porcelain --ignore-submodules        |
      | current | frontend | git push origin :current                          |
      |         | frontend | git checkout main                                 |
      |         | backend  | git rev-parse current                             |
      |         | backend  | git log main..current                             |
      | main    | frontend | git branch -D current                             |
      |         | backend  | git config --unset git-town-branch.current.parent |
      |         | backend  | git branch                                        |
      |         | backend  | git branch                                        |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git checkout other                                |
      |         | backend  | git checkout main                                 |
    And the current branch is now "main"
