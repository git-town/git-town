Feature: append a branch when the previous branch is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | current  | feature | main   | local, origin |
      | previous | feature | main   | local, origin |
    And branch "previous" is active in another worktree
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git checkout -b new      |
    And the previous Git branch is now "current"
    And no commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | new     | git checkout current |
      | current | git branch -D new    |
    And there is now no previous Git branch
    And the initial branches and lineage exist now
