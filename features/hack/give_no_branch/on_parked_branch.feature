Feature: making the current parked branch a feature branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | parked | main   | local     |
    And the current branch is "parked"
    When I run "git-town hack"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "parked" is now a feature branch
      """
    And branch "parked" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "parked" is now parked
