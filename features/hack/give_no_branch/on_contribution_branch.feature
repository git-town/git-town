Feature: making the current contribution branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the current branch is "contribution"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" is now a contribution branch
