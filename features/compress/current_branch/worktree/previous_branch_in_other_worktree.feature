Feature: compress a branch when the previous branch is active in another worktree

  Background:
    Given the feature branches "current" and "previous"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  |
      | current | local, origin | commit 1 |
      |         |               | commit 2 |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town compress"

  @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | current | git fetch --prune --tags                        |
      |         | git reset --soft main                           |
      |         | git commit -m "commit 1"                        |
      |         | git push --force-with-lease --force-if-includes |
    And the current branch is "current" and the previous branch is "previous"
    And no commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | new     | git checkout current |
      | current | git branch -D new    |
    And the current branch is now "current"
    And the initial branches and lineage exist
