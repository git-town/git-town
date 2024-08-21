Feature: park a local branch

  Background:
    Given a local Git repo
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town park feature"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now parked
      """
    And the parked branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "feature" is now a feature branch
    And there are now no parked branches
