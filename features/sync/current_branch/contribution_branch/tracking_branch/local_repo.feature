Feature: sync the current contribution branch in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the commits
      | BRANCH       | LOCATION | MESSAGE      | FILE NAME  |
      | main         | local    | main commit  | main_file  |
      | contribution | local    | local commit | local_file |
    And the current branch is "contribution"
    And the current branch is "contribution"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs no commands
    And all branches are now synchronized
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
