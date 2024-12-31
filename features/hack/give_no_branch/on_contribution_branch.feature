Feature: making the current contribution branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE   | LOCATIONS |
      | contribution | (none) | local     |
    And local Git setting "git-town.contribution-branches" is "contribution"
    And the current branch is "contribution"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" is now a feature branch
    And local Git setting "git-town.contribution-branches" is still "contribution"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" is now a contribution branch
