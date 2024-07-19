Feature: sync the current observed branch in a local repo

  Background:
    Given a local Git repo clone
    And the branch
      | NAME  | TYPE     | LOCATIONS |
      | other | observed | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | other  | local    | local commit | local_file |
    And the current branch is "other"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND       |
      | other  | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And all branches are now synchronized
    And the current branch is still "other"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | other  | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "other"
    And the initial commits exist
    And the initial branches and lineage exist
