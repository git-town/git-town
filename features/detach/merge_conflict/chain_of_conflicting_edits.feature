Feature: detaching a branch from a chain that edits the same file

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | local, origin | main commit | file      | main content |
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
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1     |
      |          | git checkout --theirs file                                     |
      |          | git add file                                                   |
      |          | git -c core.editor=true rebase --continue                      |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-4                                          |
      | branch-4 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-3 branch-2 |
      |          | git checkout --theirs file                                     |
      |          | git add file                                                   |
      |          | git -c core.editor=true rebase --continue                      |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-2                                          |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main     | local, origin | main commit | file      | main content |
      | branch-1 | local, origin | commit 1    | file      | content 1    |
      | branch-2 | local, origin | commit 2    | file      | content 2    |
      | branch-3 | local, origin | commit 2    | file      | content 2    |
      |          |               | commit 3    | file      | content 3    |
      | branch-4 | local, origin | commit 1    | file      | content 1    |
      |          |               | commit 2    | file      | content 2    |
      |          |               | commit 3    | file      | content 3    |
      |          |               | commit 4    | file      | content 4    |
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
      | branch-2 | git reset --hard {{ sha 'commit 2' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git reset --hard {{ sha 'commit 4' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial commits exist now
    And the initial lineage exists now
