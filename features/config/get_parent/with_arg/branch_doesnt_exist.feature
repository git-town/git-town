Feature: display the parent of a top-level feature branch

  Background:
    Given a Git repo clone
    When I run "git-town config get-parent zonk"

  Scenario: result
    Then it runs no commands
    And it prints no output
