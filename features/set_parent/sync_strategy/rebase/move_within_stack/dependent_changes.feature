Feature: move a branch within a stack

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 1\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                             |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                 |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And the current branch is "branch-3"
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town set-parent branch-1"

  Scenario: result
    And Git Town prints:
      """
      branch "branch-3" is now a child of "branch-1"
      """
    And Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git checkout --theirs file                                     |
      |          | git add file                                                   |
      |          | GIT_EDITOR=true git rebase --continue                          |
      |          | git push --force-with-lease --force-if-includes                |
    # TODO: the conflict above is not a phantom merge conflict
    # Because of the incorrect resolution, branch-3 still contains changes from branch-2, even though they are siblings now
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                 |
      | main     | local, origin | main commit     | file      | line 1\nline 2\nline 3                                                       |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3                                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | branch-1 |
      | branch-3 | branch-1 |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | git reset --hard {{ sha 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
    And the initial commits exist now
    And the initial branches and lineage exist now
