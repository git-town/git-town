Feature: hoisting a branch whose parent is checked out in another workspace

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-1 | local, origin | commit 1a |
      | branch-1 | local, origin | commit 1b |
    And the branches
      | NAME     | TYPE   | PARENT   | LOCATIONS     |
      | branch-2 | parked | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-2 | local, origin | commit 2a |
      | branch-2 | local, origin | commit 2b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-3 | local, origin | commit 3a |
      | branch-3 | local, origin | commit 3b |
    And branch "branch-3" is active in another worktree
    And the current branch is "branch-2"
    When I run "git-town hoist"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git fetch --prune --tags                        |
      |          | git rebase --onto main branch-1                 |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
    And Git Town prints the error:
      """
      cannot hoist branch without parent
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
