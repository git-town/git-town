Feature: with upstream repo

  Background:
    Given setting "sync-strategy" is "rebase"
    And an upstream repo
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | main    | upstream | upstream commit |
      | feature | local    | local commit    |
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git rebase origin/main      |
      |         | git fetch upstream main     |
      |         | git rebase upstream/main    |
      |         | git push                    |
      |         | git checkout feature        |
      | feature | git rebase origin/feature   |
      |         | git rebase main             |
      |         | git push --force-with-lease |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local, origin           | upstream commit |
      |         |                         | local commit    |
