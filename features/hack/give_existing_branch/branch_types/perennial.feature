Feature: hack an existing contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE      | PARENT | LOCATIONS |
      | existing | perennial | main   | local     |
    And the current branch is "main"
    When I run "git-town hack existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "existing" is a perennial branch and cannot be a feature branch
      """
    And the current branch is still "main"
    And branch "existing" still has type "perennial"
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And branch "existing" still has type "perennial"
