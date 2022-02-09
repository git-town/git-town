Feature: sync the main branch in a local repo

  Background:
    Given my repo does not have an origin
    And the current branch is "main"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs no commands
    And the current branch is still "main"
    And now the initial commits exist
