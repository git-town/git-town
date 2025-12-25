Feature: shipped the head branch of a synced stack with dependent changes that create a file while main also creates the same file

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 0\nline 1\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                             |
      | branch-1 | local, origin | branch-1 commit | file      | line 0\nline 1: branch-1 content\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | branch-2 | local, origin | branch-2 commit | file      | line 0\nline 1: branch-1 content\nline 2: branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And I add this commit to the "main" branch
      | MESSAGE                    | FILE NAME | FILE CONTENT                                |
      | independent commit on main | file      | line 0: independent content\nline 1\nline 2 |
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 0
      line 1: branch-1 content
      =======
      line 0: independent content
      line 1
      >>>>>>> {{ sha-short 'independent commit on main' }} (independent commit on main)
      line 2
      """
    When I resolve the conflict in "file" with:
      """
      line 0: independent content
      line 1: branch-1 content
      line 2
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | main     | GIT_EDITOR=true git rebase --continue                      |
      |          | git push                                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT                                                                    |
      | main     | local, origin | main commit                | file      | line 0\nline 1\nline 2                                                          |
      |          |               | branch-1 commit            | file      | line 0\nline 1: branch-1 content\nline 2                                        |
      |          |               | independent commit on main | file      | line 0: independent content\nline 1: branch-1 content\nline 2                   |
      | branch-2 | local, origin | branch-2 commit            | file      | line 0: independent content\nline 1: branch-1 content\nline 2: branch-2 content |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --abort    |
      |        | git checkout branch-2 |
    And the initial commits exist now
    And the initial branches and lineage exist now
