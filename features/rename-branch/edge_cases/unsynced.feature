Feature: rename an unsynced branch

  Background:
    Given a feature branch "old"

  Scenario: unpulled remote commits
    Given the commits
      | BRANCH | LOCATION | MESSAGE       |
      | old    | origin   | origin commit |
    And the current branch is "old"
    When I run "git-town rename-branch old new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      "old" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And the current branch is still "old"

  Scenario: unpushed local commits
    Given the commits
      | BRANCH | LOCATION | MESSAGE      |
      | old    | local    | local commit |
    And the current branch is "old"
    When I run "git-town rename-branch old new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      "old" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And the current branch is now "old"
