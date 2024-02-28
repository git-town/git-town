Feature: making multiple branches a feature branch

  Background:
    Given a contribution branch "contribution"
    And an observed branch "observed"
    And a parked branch "parked"
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
