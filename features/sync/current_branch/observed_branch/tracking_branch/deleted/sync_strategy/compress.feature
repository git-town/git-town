Feature: remove an observed branch as soon as its tracking branch is gone, even if it has unpushed commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
      | other    | observed | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE      | FILE NAME  |
      | main     | local, origin | main commit  | main_file  |
      | observed | local         | local commit | local_file |
    And the current branch is "observed"
    And origin deletes the "observed" branch
    And Git setting "git-town.sync-feature-strategy" is "compress"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git branch -D observed   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And Git Town prints:
      """
      deleted branch "observed"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | main   | git branch observed {{ sha-initial 'local commit' }} |
      |        | git checkout observed                                |
    And the initial commits exist now
    And the initial branches and lineage exist now
