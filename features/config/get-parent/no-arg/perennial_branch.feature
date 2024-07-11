Feature: display the parent of a perennial branch

  Background:
    Given the current branch is a perennial branch "perennial"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints:
      """

      """
