Feature: making a contribution branch a feature branch

  Background:
    Given the current branch is a contribution branch "contribution"
    When I run "git-town hack"

  @this
  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "feature" is now a feature branch
