Feature: non-existing branch

  Scenario:
    Given a Git repo clone
    And an uncommitted file
    When I run "git-town kill non-existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is now "main"
    And the uncommitted file still exists
