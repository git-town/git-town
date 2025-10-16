Feature: detaching a branch whose child contains merge commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-2 | local, origin | commit 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-3"
    And the current branch is "branch-2"
    And I ran "git merge branch-2 --no-edit"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-2                                          |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1     |
      |          | git push --force-with-lease --force-if-includes                |
    And this lineage exists now
      """
      main
        branch-1
          branch-3
        branch-2
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-3 | local, origin | commit 3 |
      | branch-2 | local, origin | commit 2 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial lineage exists now
    And the initial commits exist now
