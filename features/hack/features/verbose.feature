Feature: display all executed Git commands

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  Scenario: result
    When I run "git-town hack new --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva                               |
      |        | backend  | git remote                                    |
      | main   | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules                |
      | main   | frontend | git rebase origin/main                        |
      |        | backend  | git rev-list --left-right main...origin/main  |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      | main   | frontend | git branch new main                           |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git config git-town-branch.new.parent main    |
      | main   | frontend | git checkout new                              |
      |        | backend  | git show-ref --quiet refs/heads/main          |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git branch -vva                               |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 24 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules                |
      |        | backend  | git config --unset git-town-branch.new.parent |
      | new    | frontend | git checkout main                             |
      |        | backend  | git rev-parse --short HEAD                    |
      | main   | frontend | git reset --hard {{ sha 'Initial commit' }}   |
      |        | frontend | git branch -D new                             |
      |        | backend  | git show-ref --quiet refs/heads/main          |
      |        | backend  | git show-ref --quiet refs/heads/new           |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git checkout main                             |
      |        | backend  | git checkout main                             |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git branch -vva                               |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 22 shell commands.
      """
    And the current branch is now "main"
