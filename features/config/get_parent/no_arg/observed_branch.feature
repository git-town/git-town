Feature: display the parent of an observed branch

  Background:
    Given the current branch is an observed branch "observed"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints no output
