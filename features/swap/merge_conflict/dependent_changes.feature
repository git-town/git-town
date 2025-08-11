Feature: swapping a feature branch in a stack with dependent changes

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
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                   |
      | branch-2 | git fetch --prune --tags                                                                  |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                                |
      |          | git checkout --theirs file                                                                |
      |          | git add file                                                                              |
      |          | GIT_EDITOR=true git rebase --continue                                                     |
      |          | git push --force-with-lease --force-if-includes                                           |
      |          | git checkout branch-1                                                                     |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                                |
      |          | git checkout --theirs file                                                                |
      |          | git add file                                                                              |
      |          | GIT_EDITOR=true git rebase --continue                                                     |
      |          | git push --force-with-lease --force-if-includes                                           |
      |          | git checkout branch-3                                                                     |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-initial 'branch-2 commit' }} |
      |          | git checkout --theirs file                                                                |
      |          | git add file                                                                              |
      |          | GIT_EDITOR=true git rebase --continue                                                     |
      |          | git push --force-with-lease --force-if-includes                                           |
      |          | git checkout branch-2                                                                     |
    # TODO: the conflicts above are not phantom conflicts
    # In the commits below, branch-1 should contains branch-2 changes, and branch-2 should not contain branch-1 changes
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                 |
      | main     | local, origin | main commit     | file      | line 1\nline 2\nline 3                                                       |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3                                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | branch-2 |
      | branch-2 | main     |
      | branch-3 | branch-1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git checkout branch-1                           |
      | branch-1 | git reset --hard {{ sha 'branch-1 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial commits exist now
    And the initial lineage exists now
