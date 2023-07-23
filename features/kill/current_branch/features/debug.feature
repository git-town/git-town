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
    Then it runs the debug commands
      | git config -lz --local                            |
      | git config -lz --global                           |
      | git rev-parse                                     |
      | git rev-parse --show-toplevel                     |
      | git version                                       |
      | git branch -a                                     |
      | git status                                        |
      | git rev-parse --abbrev-ref HEAD                   |
      | git branch                                        |
      | git config -lz --local                            |
      | git config -lz --global                           |
      | git remote                                        |
      | git branch -r                                     |
      | git rev-parse --verify --abbrev-ref @{-1}         |
      | git status --porcelain --ignore-submodules        |
      | git rev-parse --verify --abbrev-ref @{-1}         |
      | git status --porcelain --ignore-submodules        |
      | git rev-parse current                             |
      | git log main..current                             |
      | git config --unset git-town-branch.current.parent |
      | git branch                                        |
      | git branch                                        |
      | git rev-parse --verify --abbrev-ref @{-1}         |
      | git checkout other                                |
      | git checkout main                                 |
    And the current branch is now "main"

  Scenario: undo
    Given I ran "git-town kill"
    When I run "git-town undo --debug"
    Then it runs the debug commands
      | git config -lz --local                         |
      | git config -lz --global                        |
      | git rev-parse                                  |
      | git rev-parse --show-toplevel                  |
      | git version                                    |
      | git branch -a                                  |
      | git config git-town-branch.current.parent main |
      | git status                                     |
      | git rev-parse --abbrev-ref HEAD                |
    And the current branch is now "current"
