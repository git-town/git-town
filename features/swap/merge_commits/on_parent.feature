Feature: swapping a feature branch with a branch that has merge commits

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
    And the current branch is "branch-1"
    And I ran "git merge main --message merging"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git fetch --prune --tags                        |
      |          | git rebase --onto main branch-1                 |
      |          | git checkout branch-1                           |
      | branch-1 | git rebase --onto branch-2 main                 |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git rebase --onto branch-1 branch-2             |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-1 | local, origin | commit 1    |
      | branch-2 | local, origin | commit 2    |
      | branch-3 | local, origin | commit 3    |
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
      | branch-1 | git reset --hard {{ sha 'commit 1' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-1 | local, origin | commit 1    |
      | branch-2 | local, origin | commit 2    |
      | branch-3 | local, origin | commit 3    |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | branch-1 |
      | branch-3 | branch-2 |
