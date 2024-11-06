Feature: non-existing branch

  Scenario:
    Given a Git repo with origin
    And an uncommitted file
    When I run "git-town delete non-existing"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is now "main"
    And the uncommitted file still exists
