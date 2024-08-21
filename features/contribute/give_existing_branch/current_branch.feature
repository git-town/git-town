Feature: make the current contribution branch a contribution branch

  Background:
    Given a local Git repo
    And the branch
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution | main   | local     |
    And the current branch is "contribution"
    And an uncommitted file
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "contribution" is already a contribution branch
      """
    And the contribution branches are still "contribution"
    And the current branch is still "contribution"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "contribution"
    And the contribution branches are still "contribution"
    And the uncommitted file still exists
