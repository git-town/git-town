Feature: remove the observed branch as soon as the tracking branch is gone, even if it has unpushed commits

  Background:
    Given the current branch is an observed branch "observed"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE      | FILE NAME  |
      | main     | local, origin | main commit  | main_file  |
      | observed | local         | local commit | local_file |
    And origin deletes the "observed" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git branch -D observed   |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And it prints:
      """
      deleted branch "observed"
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                 |
      | main   | git branch observed {{ sha-before-run 'local commit' }} |
      |        | git checkout observed                                   |
    And the current branch is now "observed"
    And the initial commits exist
    And the initial branches and lineage exist
