Feature: append a branch when the previous branch is active in another worktree

  Background:
    Given the feature branches "current" and "previous"
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | current | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout current                    |
      | current | git merge --no-edit --ff origin/current |
      |         | git merge --no-edit --ff main           |
      |         | git checkout -b new                     |
    And the current branch is "new" and the previous branch is "current"
    And no commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | new     | git checkout current |
      | current | git branch -D new    |
    And the current branch is now "current"
    And the initial branches and lineage exist
