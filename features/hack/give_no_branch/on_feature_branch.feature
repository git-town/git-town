Feature: cannot make the current feature branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "feature" is already a feature branch
      """
    And branch "feature" is still a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" is still a feature branch
