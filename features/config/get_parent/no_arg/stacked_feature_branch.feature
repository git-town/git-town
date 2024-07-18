Feature: display the parent of a stacked feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      parent
      """
