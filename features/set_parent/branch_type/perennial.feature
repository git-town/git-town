Feature: cannot set parent of perennial branches

  Background:
    Given a Git repo with origin

  Scenario: on main branch
    When I run "git-town set-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      the branch main is not a feature branch. Only feature branches can have parent branches
      """
    And the initial lineage exists now

  Scenario: on perennial branch
    Given the branches
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the current branch is "qa"
    When I run "git-town set-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      the branch qa is not a feature branch. Only feature branches can have parent branches
      """
    And the initial lineage exists now
