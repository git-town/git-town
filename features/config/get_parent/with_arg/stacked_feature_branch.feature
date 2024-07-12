Feature: display the parent of a stacked feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    When I run "git-town config get-parent child"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      parent
      """
