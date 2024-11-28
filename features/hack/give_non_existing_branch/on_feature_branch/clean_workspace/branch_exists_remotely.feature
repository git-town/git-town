Feature: already existing remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    When I run "git-town hack existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is already a branch "existing" toRefId the "origin" remote
      """
    And the current branch is still "main"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
