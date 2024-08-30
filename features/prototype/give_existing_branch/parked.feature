Feature: prototype another parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "main"
    When I run "git-town prototype parked"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "parked" is now a prototype branch
      """
    And the prototype branches are now "parked"
    And the parked branches are still "parked"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the parked branches are still "parked"
