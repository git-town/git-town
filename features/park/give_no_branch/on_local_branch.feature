Feature: park the current local branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now parked
      """
    And the current branch is still "feature"
    And the parked branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And there are now no parked branches
