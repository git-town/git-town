Feature: cannot set parent of perennial branches

  Background:
    Given a Git repo clone

  Scenario: on main branch
    When I run "git-town set-parent"
    Then it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And it runs no commands
    And the initial lineage exists
    And the current branch is still "main"

  Scenario: on perennial branch
    Given the branches
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the current branch is "qa"
    When I run "git-town set-parent"
    Then it prints the error:
      """
      the branch "qa" is not a feature branch. Only feature branches can have parent branches
      """
    And it runs no commands
    And the initial lineage exists
    And the current branch is still "qa"
