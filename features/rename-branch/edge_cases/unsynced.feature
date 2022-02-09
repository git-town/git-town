Feature: rename an unsynced branch

  Background:
    Given a feature branch "old"

  Scenario: unpulled remote commits
    Given the commits
      | BRANCH | LOCATION | MESSAGE       |
      | old    | origin   | origin commit |
    And I am on the "old" branch
    When I run "git-town rename-branch old new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      "old" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am still on the "old" branch

  Scenario: unpushed local commits
    Given the commits
      | BRANCH | LOCATION | MESSAGE      |
      | old    | local    | local commit |
    And I am on the "old" branch
    When I run "git-town rename-branch old new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      "old" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am now on the "old" branch
