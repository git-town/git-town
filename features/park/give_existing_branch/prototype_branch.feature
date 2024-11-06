Feature: park another prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town park prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "prototype" is now parked
      """
    And the prototype branches are now "prototype"
    And the parked branches are now "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the prototype branches are now "prototype"
    And there are now no parked branches
    And the initial branches and lineage exist now
