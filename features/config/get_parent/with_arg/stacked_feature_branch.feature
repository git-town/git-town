Feature: display the parent of a stacked feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    When I run "git-town config get-parent child"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      parent
      """
