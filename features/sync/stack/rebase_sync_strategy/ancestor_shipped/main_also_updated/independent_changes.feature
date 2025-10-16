Feature: shipped the head branch of a synced stack with inddependent changes that create a file while main also creates the same file with independent changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT               |
      | main   | local, origin | main commit | file      | line 0\n\nline 1\n\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                 |
      | branch-1 | local, origin | branch-1 commit | file      | line 0\n\nline 1: branch-1 content\n\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 0\n\nline 1: branch-1 content\n\nline 2: branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And I add this commit to the "main" branch
      | MESSAGE                    | FILE NAME | FILE CONTENT                                    |
      | independent commit on main | file      | line 0: independent content\n\nline 1\n\nline 2 |
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git push                                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT                                                                        |
      | main     | local, origin | main commit                | file      | line 0\n\nline 1\n\nline 2                                                          |
      |          |               | branch-1 commit            | file      | line 0\n\nline 1: branch-1 content\n\nline 2                                        |
      |          |               | independent commit on main | file      | line 0: independent content\n\nline 1: branch-1 content\n\nline 2                   |
      | branch-2 | local, origin | branch-2 commit            | file      | line 0: independent content\n\nline 1: branch-1 content\n\nline 2: branch-2 content |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                 |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}            |
      |          | git push --force-with-lease --force-if-includes         |
      |          | git branch branch-1 {{ sha-initial 'branch-1 commit' }} |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT                                                      |
      | main     | local, origin | main commit                | file      | line 0\n\nline 1\n\nline 2                                        |
      |          |               | branch-1 commit            | file      | line 0\n\nline 1: branch-1 content\n\nline 2                      |
      |          |               | independent commit on main | file      | line 0: independent content\n\nline 1: branch-1 content\n\nline 2 |
      | branch-1 | local         | branch-1 commit            | file      | line 0\n\nline 1: branch-1 content\n\nline 2                      |
      | branch-2 | local, origin | branch-2 commit            | file      | line 0\n\nline 1: branch-1 content\n\nline 2: branch-2 content    |
      |          | origin        | branch-1 commit            | file      | line 0\n\nline 1: branch-1 content\n\nline 2                      |
