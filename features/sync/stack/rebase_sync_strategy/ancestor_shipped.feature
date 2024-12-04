Feature: shipped the head branch of a synced stack with dependent changes

  Background:
    Given a Git repo with origin
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
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-3 | local, origin | commit 3 | file      | content 3    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-4 | feature | branch-3 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-4 | local, origin | commit 4 | file      | content 4    |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "branch-4"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And origin ships the "branch-2" branch using the "squash-merge" ship-strategy and resolves the merge conflict in "file" with "content 2" and commits as "commit 2"
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-4 | git fetch --prune --tags                        |
      |          | git checkout main                               |
      | main     | git rebase origin/main --no-update-refs         |
      |          | git branch -D branch-1                          |
      |          | git branch -D branch-2                          |
      |          | git checkout branch-3                           |
      | branch-3 | git rebase main --no-update-refs                |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git rebase branch-3 --no-update-refs            |
      |          | git push --force-with-lease --force-if-includes |
    And the current branch is now "branch-4"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | main     | local, origin | commit 1 | file      | content 1    |
      |          |               | commit 2 | file      | content 2    |
      | branch-3 | local, origin | commit 3 | file      | content 3    |
      | branch-4 | local, origin | commit 4 | file      | content 4    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git checkout main                                    |
      | main   | git reset --hard {{ sha 'initial commit' }}          |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git checkout beta                                    |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | main   | origin        | alpha commit | file      | alpha content |
      | alpha  | local         | alpha commit | file      | alpha content |
      | beta   | local, origin | beta commit  | file      | beta content  |
      |        | origin        | alpha commit | file      | alpha content |
    And the initial branches and lineage exist now
