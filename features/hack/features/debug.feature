Feature: display debug statistics

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  @debug @this
  Scenario: result
    When I run "git-town hack new --debug"
    Then it runs the debug commands
      | git config -lz --local                       |
      | git config -lz --global                      |
      | git rev-parse                                |
      | git rev-parse --show-toplevel                |
      | git version                                  |
      | git branch -a                                |
      | git remote                                   |
      | git status                                   |
      | git rev-parse --abbrev-ref HEAD              |
      | git branch -a                                |
      | git branch -r                                |
      | git rev-parse --verify --abbrev-ref @{-1}    |
      | git status --porcelain --ignore-submodules   |
      | git rev-parse HEAD                           |
      | git rev-list --left-right main...origin/main |
      | git config git-town-branch.new.parent main   |
      | git branch                                   |
      | git branch                                   |
      | git rev-parse --verify --abbrev-ref @{-1}    |

    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git town undo --debug"
    Then it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is now "main"
