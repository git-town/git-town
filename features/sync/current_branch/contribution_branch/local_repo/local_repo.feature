Feature: sync the current contribution branch in a local repo

  Background:
    Given a local Git repo clone
    And the branch
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the current branch is "contribution"
    And the commits
      | BRANCH       | LOCATION | MESSAGE      | FILE NAME  |
      | main         | local    | main commit  | main_file  |
      | contribution | local    | local commit | local_file |
    And the current branch is "contribution"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND       |
      | contribution | git add -A    |
      |              | git stash     |
      |              | git stash pop |
    And all branches are now synchronized
    And the current branch is still "contribution"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND       |
      | contribution | git add -A    |
      |              | git stash     |
      |              | git stash pop |
    And the current branch is still "contribution"
    And the initial commits exist
    And the initial branches and lineage exist
