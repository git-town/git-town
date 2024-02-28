Feature: making the current observed branch a feature branch

  Background:
    Given the current branch is an observed branch "observed"
    When I run "git-town hack"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "observed" is now a feature branch
      """
    And branch "observed" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "observed" is now observed
