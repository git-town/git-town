Feature: missing configuration

  Background: running unconfigured
    Given Git Town is not configured
    When I run "git-town hack feature" and enter into the dialog:
      | DIALOG                  | KEYS  |
      | main development branch | enter |

  Scenario: result
    And it runs the commands
      | BRANCH | COMMAND                      |
      | main   | git fetch --prune --tags     |
      |        | git rebase origin/main       |
      |        | git checkout -b feature main |
    And the main branch is now "main"
    And the current branch is now "feature"
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -D feature |
    And the current branch is now "main"
    And no lineage exists now
