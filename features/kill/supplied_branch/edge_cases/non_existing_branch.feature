Feature: non-existing branch

  Scenario:
    Given the current branch is "main"
    And my workspace has an uncommitted file
    When I run "git-town kill non-existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch named "non-existing"
      """
    And the current branch is now "main"
    And my workspace still contains my uncommitted file
