Feature: park the current prototype branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the current branch is "prototype"
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "prototype" is now parked
      """
    And the parked branches are now "prototype"
    And the prototype branches are still "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are now no parked branches
    And the prototype branches are still "prototype"
