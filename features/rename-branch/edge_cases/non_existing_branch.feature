Feature: branch to rename does not exist

  Scenario: unknown branch
    Given my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
    And I am on the "main" branch
    When I run "git-town rename-branch non-existing-feature renamed-feature"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch named "non-existing-feature"
      """
    And I am still on the "main" branch
