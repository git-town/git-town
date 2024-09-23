Feature: make another prototype branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    When I run "git-town contribute prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "prototype" is now a contribution branch
      """
    And there are now no prototype branches
    And the contribution branches are now "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the prototype branches are now "prototype"
    And there are now no parked branches
    And the initial branches and lineage exist now
