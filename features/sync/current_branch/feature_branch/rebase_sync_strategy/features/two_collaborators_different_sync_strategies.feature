Feature: collaborative feature branch syncing

  Scenario: I sync via "rebase" and my coworker sync via "merge"
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And a coworker clones the repository
    And the current branch is a feature branch "feature"
    And the coworker fetches updates
    And the coworker sets the parent branch of "feature" as "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | feature | local    | my commit       |
      |         | coworker | coworker commit |
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git rebase origin/main      |
      |         | git checkout feature        |
      | feature | git rebase origin/feature   |
      |         | git rebase main             |
      |         | git push --force-with-lease |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         |
      | feature | local, origin | my commit       |
      |         | coworker      | coworker commit |
    And all branches are now synchronized
