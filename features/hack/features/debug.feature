Feature: display debug statistics

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  Scenario: result
    When I run "git-town hack new --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                      |
      |        | backend  | git version                                  |
      |        | backend  | git config -lz --local                       |
      |        | backend  | git config -lz --global                      |
      |        | backend  | git rev-parse --show-toplevel                |
      |        | backend  | git remote                                   |
      |        | backend  | git status                                   |
      |        | backend  | git rev-parse --abbrev-ref HEAD              |
      | main   | frontend | git fetch --prune --tags                     |
      |        | backend  | git branch -vva                              |
      |        | backend  | git branch -a                                |
      |        | backend  | git rev-parse --show-toplevel                |
      |        | backend  | git branch -r                                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}    |
      |        | backend  | git status --porcelain --ignore-submodules   |
      |        | backend  | git rev-parse HEAD                           |
      | main   | frontend | git rebase origin/main                       |
      |        | backend  | git rev-list --left-right main...origin/main |
      | main   | frontend | git branch new main                          |
      |        | backend  | git config git-town-branch.new.parent main   |
      | main   | frontend | git checkout new                             |
      |        | backend  | git branch                                   |
      |        | backend  | git branch                                   |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}    |
    And it prints:
      """
      Ran 23 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git town undo --debug"
    Then it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is now "main"
