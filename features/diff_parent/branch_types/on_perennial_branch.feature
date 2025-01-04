Feature: on perennial branch

  Scenario: on main branch
    Given a Git repo with origin
    And the current branch is "main"
    When I run "git-town diff-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: on perennial branch
    Given a Git repo with origin
    And the branches
      | NAME | TYPE      | LOCATIONS |
      | qa   | perennial | local     |
    And the current branch is "qa"
    When I run "git-town diff-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you can only diff-parent feature branches
      """
