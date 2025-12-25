Feature: sync the main branch in a local repo

  Background:
    Given a local Git repo
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
    And the current branch is "main"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs no commands
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
