Feature: sync the current perennial branch in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME       | TYPE      | LOCATIONS |
      | production | perennial | local     |
      | qa         | perennial | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | qa     | local    | local commit | local_file |
    And the current branch is "qa"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And all branches are now synchronized
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | main commit  |
      | qa     | local    | local commit |
    And the initial branches and lineage exist now
