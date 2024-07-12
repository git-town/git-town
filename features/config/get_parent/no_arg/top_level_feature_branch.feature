Feature: display the parent of a top-level feature branch

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      main
      """
