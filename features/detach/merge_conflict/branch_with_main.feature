Feature: detaching a branch that conflicts with the main branch

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
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git checkout --theirs file                                 |
      |          | git add file                                               |
      |          | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease --force-if-includes            |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main     | local, origin | main commit | file      | main content |
      | branch-1 | local, origin | commit 1    | file      | content 1    |
      | branch-2 | local, origin | commit 2    | file      | content 2    |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-1 | main   |
      | branch-2 | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}           |
      |          | git push --force-with-lease --force-if-includes |
    And the initial commits exist now
    And the initial lineage exists now
