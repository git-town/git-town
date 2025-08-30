Feature: prototype another prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town prototype prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "prototype" is already a prototype branch
      """
    And branch "prototype" still has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "prototype" still has type "prototype"
    And the initial branches and lineage exist now
