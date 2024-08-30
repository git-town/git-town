Feature: park the current parked branch

  Background:
    Given a local Git repo
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | parked | main   | local     |
    And the current branch is "parked"
    When I run "git-town park parked"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "parked" is already parked
      """
    And the parked branches are still "parked"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the parked branches are still "parked"
