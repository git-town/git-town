Feature: display all executed Git commands

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |

  Scenario: result
    When I run "git-town append new --verbose"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                              |
      |          | backend  | git version                                          |
      |          | backend  | git config -lz --global                              |
      |          | backend  | git config -lz --local                               |
      |          | backend  | git rev-parse --show-toplevel                        |
      |          | backend  | git stash list                                       |
      |          | backend  | git branch -vva                                      |
      |          | backend  | git remote                                           |
      | existing | frontend | git fetch --prune --tags                             |
      |          | backend  | git branch -vva                                      |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}            |
      |          | backend  | git status --ignore-submodules                       |
      | existing | frontend | git checkout main                                    |
      | main     | frontend | git rebase origin/main                               |
      |          | backend  | git rev-list --left-right main...origin/main         |
      | main     | frontend | git checkout existing                                |
      | existing | frontend | git merge --no-edit origin/existing                  |
      |          | frontend | git merge --no-edit main                             |
      |          | backend  | git rev-list --left-right existing...origin/existing |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing    |
      | existing | frontend | git branch new existing                              |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing    |
      |          | backend  | git config git-town-branch.new.parent existing       |
      | existing | frontend | git checkout new                                     |
      |          | backend  | git show-ref --quiet refs/heads/existing             |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}            |
      |          | backend  | git config -lz --global                              |
      |          | backend  | git config -lz --local                               |
      |          | backend  | git branch -vva                                      |
      |          | backend  | git stash list                                       |
    And it prints:
      """
      Ran 29 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --verbose"
    Then it prints:
      """
      Ran 20 shell commands.
      """
    And the current branch is now "existing"
