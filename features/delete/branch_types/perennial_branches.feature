Feature: does not delete perennial branches

  Scenario: main branch
    Given a Git repo with origin
    When I run "git-town delete"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      you cannot delete the main branch
      """

  Scenario: perennial branch
    Given a Git repo with origin
    And the branches
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the current branch is "qa"
    When I run "git-town delete"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
    And Git Town prints the error:
      """
      you cannot delete perennial branches
      """
