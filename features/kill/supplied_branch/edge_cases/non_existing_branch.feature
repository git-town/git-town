Feature: non-existing branch

  Scenario:
    Given I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town kill non-existing-feature"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch named "non-existing-feature"
      """
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
