Feature: display the parent of an observed branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    And the current branch is "observed"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints no output
