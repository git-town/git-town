Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |

  Scenario: result
    When I run "git-town append new --debug"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                              |
      |          | backend  | git version                                          |
      |          | backend  | git config -lz --local                               |
      |          | backend  | git config -lz --global                              |
      |          | backend  | git rev-parse --show-toplevel                        |
      |          | backend  | git rev-parse --show-toplevel                        |
      |          | backend  | git remote                                           |
      |          | backend  | git status                                           |
      |          | backend  | git rev-parse --abbrev-ref HEAD                      |
      | existing | frontend | git fetch --prune --tags                             |
      |          | backend  | git branch -vva                                      |
      |          | backend  | git branch -a                                        |
      |          | backend  | git branch -r                                        |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}            |
      |          | backend  | git status --porcelain --ignore-submodules           |
      | existing | frontend | git checkout main                                    |
      |          | backend  | git rev-parse HEAD                                   |
      | main     | frontend | git rebase origin/main                               |
      |          | backend  | git rev-list --left-right main...origin/main         |
      | main     | frontend | git checkout existing                                |
      |          | backend  | git rev-parse HEAD                                   |
      | existing | frontend | git merge --no-edit origin/existing                  |
      |          | backend  | git rev-parse HEAD                                   |
      | existing | frontend | git merge --no-edit main                             |
      |          | backend  | git rev-list --left-right existing...origin/existing |
      | existing | frontend | git branch new existing                              |
      |          | backend  | git config git-town-branch.new.parent existing       |
      | existing | frontend | git checkout new                                     |
      |          | backend  | git branch                                           |
      |          | backend  | git branch                                           |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}            |
    And it prints:
      """
      Ran 30 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 14 shell commands.
      """
    And the current branch is now "existing"
