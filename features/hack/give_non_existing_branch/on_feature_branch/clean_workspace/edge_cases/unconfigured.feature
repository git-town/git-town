@messyoutput
Feature: missing configuration

  Background: running unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town hack feature" and enter into the dialog:
      | DIALOG      | KEYS  |
      | main branch | enter |

  Scenario: result
    And Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git rebase origin/main --no-update-refs |
      |        | git checkout -b feature                 |
    And the main branch is now "main"
    And the current branch is now "feature"
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -D feature |
    And the current branch is now "main"
    And no lineage exists now
