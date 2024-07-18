Feature: branch does not exist

  Scenario:
    Given a Git repo clone
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And the current branch is "main"
    When I run "git-town rename-branch non-existing new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is still "main"
