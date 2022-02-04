Feature: sync the main branch in a local repo

  Background:
    Given my repo does not have a remote origin
    And I am on the "main" branch
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs no commands
    And I am still on the "main" branch
    And my repo is left with my original commits
