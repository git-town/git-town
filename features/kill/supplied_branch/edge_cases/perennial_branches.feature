Feature: does not kill perennial branches

  Scenario: main branch
    Given the current branch is a feature branch "feature"
    When I run "git-town kill main"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot kill the main branch
      """
    And the current branch is still "feature"
    And the initial branches and lineage exist

  Scenario: perennial branch
    Given a perennial branch "qa"
    And the current branch is "main"
    When I run "git-town kill qa"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot kill perennial branches
      """
