Feature: swapping a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | parked  | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-1 | local, origin | commit 1a   |
      | branch-1 | local, origin | commit 1b   |
      | branch-2 | local, origin | commit 2a   |
      | branch-2 | local, origin | commit 2b   |
      | branch-3 | local, origin | commit 3a   |
      | branch-3 | local, origin | commit 3b   |
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                             |
      | branch-2 | git fetch --prune --tags                                                            |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                          |
      |          | git push --force-with-lease --force-if-includes                                     |
      |          | git checkout branch-1                                                               |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                          |
      |          | git push --force-with-lease --force-if-includes                                     |
      |          | git checkout branch-3                                                               |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-initial 'commit 2b' }} |
      |          | git push --force-with-lease --force-if-includes                                     |
      |          | git checkout branch-2                                                               |
    And this lineage exists now
      """
      main
        branch-2
          branch-1
            branch-3
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-2 | local, origin | commit 2a   |
      |          |               | commit 2b   |
      | branch-1 | local, origin | commit 1a   |
      |          |               | commit 1b   |
      | branch-3 | local, origin | commit 3a   |
      |          |               | commit 3b   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git checkout branch-1                           |
      | branch-1 | git reset --hard {{ sha 'commit 1b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial lineage exists now
    And the initial commits exist now
