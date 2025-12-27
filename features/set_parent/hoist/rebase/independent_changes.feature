Feature: remove a branch and all its children from a stack with independent changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT               |
      | main   | local, origin | main commit | file      | line 1\n\nline 2\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                 |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\n\nline 2\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\n\nline 2: branch-2 changes\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                     |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\n\nline 2: branch-2 changes\n\nline 3: branch-3 changes |
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-2"
    When I run "git-town set-parent main"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1     |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-2 branch-1 |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout branch-2                                          |
    And Git Town prints:
      """
      branch branch-2 is now a child of main
      """
    And this lineage exists now
      """
      main
        branch-1
        branch-2
          branch-3
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                   |
      | main     | local, origin | main commit     | file      | line 1\n\nline 2\n\nline 3                                     |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\n\nline 2\n\nline 3                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1\n\nline 2: branch-2 changes\n\nline 3                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 1\n\nline 2: branch-2 changes\n\nline 3: branch-3 changes |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial branches and lineage exist now
    And the initial commits exist now
