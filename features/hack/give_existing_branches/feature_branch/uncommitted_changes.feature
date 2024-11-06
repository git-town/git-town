Feature: already existing local branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And an uncommitted file
    When I run "git-town hack existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "existing" is already a feature branch
      """
    And the current branch is still "main"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
