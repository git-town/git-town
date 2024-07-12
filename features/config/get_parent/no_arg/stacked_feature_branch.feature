Feature: display the parent of a stacked feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      parent
      """
