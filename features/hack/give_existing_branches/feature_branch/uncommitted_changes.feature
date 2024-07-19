Feature: already existing local branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And an uncommitted file
    When I run "git-town hack existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "existing" is already a feature branch
      """
    And the current branch is still "main"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the initial commits exist
    And the initial branches and lineage exist
