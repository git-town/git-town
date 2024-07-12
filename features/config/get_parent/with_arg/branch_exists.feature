Feature: display the parent of a top-level feature branch

  Background:
    Given a feature branch "feature"
    When I run "git-town config get-parent feature"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      main
      """
