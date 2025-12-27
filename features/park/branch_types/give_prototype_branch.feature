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
      branch prototype is now parked
      """
    And branch "prototype" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And branch "prototype" now has type "prototype"
