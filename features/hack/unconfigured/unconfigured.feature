@messyoutput
Feature: missing configuration

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town hack feature" and enter into the dialog:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches |       |
      | origin hostname    | enter |
      | forge type         | enter |
      | enter all          | enter |
      | config storage     | enter |

  Scenario: result
    And Git Town runs the commands
      | BRANCH | COMMAND                              |
      | main   | git fetch --prune --tags             |
      |        | git config git-town.main-branch main |
      |        | git checkout -b feature              |
    And this lineage exists now
      """
      main
        feature
      """
    And the main branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -D feature |
    And no lineage exists now
