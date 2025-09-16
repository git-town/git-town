Feature: disable pushing through the config file

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And the committed configuration file:
      """
      [branches]
      main = "main"

      [sync]
      push-branches = false
      """
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And the current branch is "branch-2"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                      |
      | branch-2 | git fetch --prune --tags     |
      |          | git checkout branch-1        |
      | branch-1 | git checkout branch-2        |
      | branch-2 | git checkout -b new branch-1 |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout branch-2 |
      | branch-2 | git branch -D new     |
    And the initial commits exist now
    And the initial lineage exists now
    And the initial tags exist now
