Feature: making multiple branches a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution |        | local     |
      | observed     | observed     |        | local     |
      | parked       | parked       | main   | local     |
    When I run "git-town hack contribution observed parked"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" is now a feature branch
    And it prints:
      """
      branch "observed" is now a feature branch
      """
    And branch "observed" is now a feature branch
    And it prints:
      """
      branch "parked" is now a feature branch
      """
    And branch "parked" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "contribution" is now a contribution branch
    And branch "observed" is now observed
    And branch "parked" is now parked
