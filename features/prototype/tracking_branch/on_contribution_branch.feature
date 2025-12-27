Feature: prototype the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the current branch is "contribution"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch contribution is now a prototype branch
      """
    And branch "contribution" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And branch "contribution" now has type "contribution"
