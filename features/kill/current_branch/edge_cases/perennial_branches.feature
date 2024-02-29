Feature: does not kill perennial branches

  Scenario: main branch
    Given the current branch is "main"
    When I run "git-town kill"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot kill the main branch
      """
    And the current branch is still "main"

  Scenario: perennial branch
    Given the current branch is a perennial branch "qa"
    When I run "git-town kill"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot kill perennial branches
      """
    And the current branch is still "qa"
