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
      |          | backend  | git status --long --ignore-submodules                |
      |          | backend  | git remote                                           |
      |          | backend  | git rev-parse --abbrev-ref HEAD                      |
      | existing | frontend | git fetch --prune --tags                             |
      |          | backend  | git stash list                                       |
      |          | backend  | git branch -vva --sort=refname                       |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}            |
      | existing | frontend | git checkout main                                    |
      | main     | frontend | git rebase origin/main                               |
      |          | backend  | git rev-list --left-right main...origin/main         |
      | main     | frontend | git checkout existing                                |
      | existing | frontend | git merge --no-edit --ff origin/existing             |
      |          | frontend | git merge --no-edit --ff main                        |
      |          | backend  | git rev-list --left-right existing...origin/existing |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing    |
      | existing | frontend | git checkout -b new                                  |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing    |
      |          | backend  | git config git-town-branch.new.parent existing       |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing    |
      |          | backend  | git branch -vva --sort=refname                       |
      |          | backend  | git config -lz --global                              |
      |          | backend  | git config -lz --local                               |
      |          | backend  | git stash list                                       |
    And it prints:
      """
      Ran 27 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                       |
      |          | backend  | git version                                   |
      |          | backend  | git config -lz --global                       |
      |          | backend  | git config -lz --local                        |
      |          | backend  | git rev-parse --show-toplevel                 |
      |          | backend  | git status --long --ignore-submodules         |
      |          | backend  | git stash list                                |
      |          | backend  | git branch -vva --sort=refname                |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |          | backend  | git remote get-url origin                     |
      | new      | frontend | git checkout existing                         |
      | existing | frontend | git branch -D new                             |
      |          | backend  | git config --unset git-town-branch.new.parent |
    And it prints:
      """
      Ran 12 shell commands.
      """
    And the current branch is still "existing"
    And the initial commits exist
    And the initial branches and lineage exist
