Feature: already existing local branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    When I run "git-town feature existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch existing is already a feature branch
      """
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
