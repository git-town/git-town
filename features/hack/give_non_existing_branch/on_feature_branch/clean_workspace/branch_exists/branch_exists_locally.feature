Feature: already existing local branch

  Background:
    Given a local feature branch "existing"
    When I run "git-town hack existing"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      branch "existing" is already a feature branch
      """
    And the current branch is still "main"
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the initial commits exist
    And the initial branches and lineage exist
