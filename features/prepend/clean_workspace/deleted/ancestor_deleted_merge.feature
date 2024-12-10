Feature: prepend a branch to a branch that was shipped at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
    And origin deletes the "branch-1" branch
    And Git Town setting "sync-feature-strategy" is "merge"
    And the current branch is "branch-2"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git fetch --prune --tags                 |
      |          | git checkout main                        |
      | main     | git rebase origin/main --no-update-refs  |
      |          | git branch -D branch-1                   |
      |          | git checkout branch-2                    |
      | branch-2 | git merge --no-edit --ff main            |
      |          | git merge --no-edit --ff origin/branch-2 |
      |          | git checkout -b new main                 |
    And Git Town prints:
      """
      deleted branch "branch-1"
      """
    And Git Town prints:
      """
      branch "branch-2" is now a child of "new"
      """
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES            |
      | local      | main, branch-2, new |
      | origin     | main, branch-2      |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-2 | new    |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | new      | git branch branch-1 {{ sha 'commit 1' }} |
      |          | git checkout branch-2                    |
      | branch-2 | git branch -D new                        |
    And the current branch is now "branch-2"
    And the initial branches and lineage exist now
