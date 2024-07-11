Feature: display the parent of a contribution branch

  Background:
    Given the current branch is a contribution branch "contribution"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints:
      """

      """
