Feature: display the parent of a top-level feature branch

  Background:
    Given a Git repo with origin
    When I run "git-town config get-parent zonk"

  Scenario: result
    Then Git Town runs no commands
    And it prints no output
