Feature: detaching a branch from a stack with independent changes

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                                 |
      | main   | local, origin | main commit | file      | line 0: main content\n\nline 1\n\nline 2\n\nline 3\n\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                   |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3\n\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2: branch-2 content\n\nline 3\n\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                                       |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2: branch-2 content\n\nline 3: branch-3 content\n\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-4 | feature | branch-3 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                                                 |
      | branch-4 | local, origin | branch-4 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3: branch-3 content\nline 4: branch-4 content |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-4                                          |
      | branch-4 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-3 branch-2 |
      |          | git checkout --theirs file                                     |
      |          | git add file                                                   |
      |          | GIT_EDITOR=true git rebase --continue                          |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-2                                          |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1     |
      |          | git push --force-with-lease --force-if-includes                |
    # TODO: branch-4 below contains changes from branch-2 but shouldn't
    # This happens because Git Town wrongfully assumes a phantom merge conflict when removing branch-2 from branch-4.
    # It works correctly for branch-3.
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                                                 |
      | main     | local, origin | main commit     | file      | line 0: main content\n\nline 1\n\nline 2\n\nline 3\n\nline 4                                                                 |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3\n\nline 4                                               |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\n\nline 1\n\nline 2: branch-2 content\n\nline 3\n\nline 4                                               |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3: branch-3 content\n\nline 4                             |
      | branch-4 | local, origin | branch-4 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3: branch-3 content\nline 4: branch-4 content |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | main     |
      | branch-3 | branch-1 |
      | branch-4 | branch-3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git reset --hard {{ sha 'branch-4 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial commits exist now
    And the initial lineage exists now
