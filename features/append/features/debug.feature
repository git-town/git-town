Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |

  Scenario: result
    When I run "git-town append new --debug"
    Then it runs the debug commands
      | git config -lz --local                               |
      | git config -lz --global                              |
      | git rev-parse                                        |
      | git rev-parse --show-toplevel                        |
      | git version                                          |
      | git branch -a                                        |
      | git status                                           |
      | git rev-parse --abbrev-ref HEAD                      |
      | git remote                                           |
      | git branch -a                                        |
      | git branch -r                                        |
      | git rev-parse --verify --abbrev-ref @{-1}            |
      | git status --porcelain --ignore-submodules           |
      | git rev-parse HEAD                                   |
      | git rev-list --left-right main...origin/main         |
      | git rev-parse HEAD                                   |
      | git rev-parse HEAD                                   |
      | git rev-list --left-right existing...origin/existing |
      | git config git-town-branch.new.parent existing       |
      | git branch                                           |
      | git branch                                           |
      | git rev-parse --verify --abbrev-ref @{-1}            |
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 15 shell commands.
      """
    And the current branch is now "existing"
