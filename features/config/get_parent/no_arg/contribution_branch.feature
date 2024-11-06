Feature: display the parent of a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
    And the current branch is "contribution"
    When I run "git-town config get-parent"

  Scenario: result
    Then Git Town runs no commands
    And it prints no output
