Feature: non-existing branch

  @debug @this
  Scenario:
    Given the current branch is "main"
    And an uncommitted file
    When I run "git-town kill non-existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      branch "non-existing" does not exist
      """
    And the current branch is now "main"
    And the uncommitted file still exists
