Feature: stack that changes the same file in multiple commits per branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1 | file      | content 1    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-2 | local, origin | commit 2 | file      | content 2    |
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-2 | git fetch --prune --tags |
      |          | git checkout branch-1    |
      | branch-1 | git checkout branch-2    |
    And no rebase is now in progress
    And all branches are now synchronized
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
