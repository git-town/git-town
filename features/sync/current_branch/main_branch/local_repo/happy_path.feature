Feature: sync the main branch in a local repo

  Background:
    Given a Git repo clone
    And my repo does not have an origin
    And the current branch is "main"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs no commands
    And the current branch is still "main"
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the initial commits exist
    And the initial branches and lineage exist
