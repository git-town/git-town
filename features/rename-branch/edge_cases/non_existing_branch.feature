Feature: branch does not exist

  Scenario:
    Given my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And I am on the "main" branch
    When I run "git-town rename-branch non-existing new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch named "non-existing"
      """
    And I am still on the "main" branch
