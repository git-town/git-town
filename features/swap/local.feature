Feature: detaching an omni-branch

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      | branch-1 | local    | commit 1b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-2 | feature | branch-1 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-2 | local    | commit 2a |
      | branch-2 | local    | commit 2b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-3 | feature | branch-2 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-3 | local    | commit 3a |
      | branch-3 | local    | commit 3b |
    And the current branch is "branch-2"
    # # branch-2 goes on top of main, removing the branch-1 commits
    # When I run "git rebase --onto main branch-1"
    # # branch-1 goes on top of the new branch-2, removing main commits
    # And I run "git checkout branch-1"
    # And I run "git rebase --onto branch-2 main"
    # # branch-3 goes on top of the new branch-1, removing the old branch-2 commits
    # And I run "git checkout branch-3"
    # And I run "git rebase --onto branch-1 branch-2"
    # And I run "git checkout branch-2"
    # # update lineage
    # And I run "git config git-town-branch.branch-2.parent main"
    # And I run "git config git-town-branch.branch-1.parent branch-2"
    # And I run "git config git-town-branch.branch-3.parent branch-1"
    When I run "git-town swap"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE     |
      | main     | local    | main commit |
      | branch-1 | local    | commit 1a   |
      |          |          | commit 1b   |
      | branch-2 | local    | commit 2a   |
      |          |          | commit 2b   |
      | branch-3 | local    | commit 3a   |
      |          |          | commit 3b   |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | branch-2 |
      | branch-2 | main     |
      | branch-3 | branch-1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git reset --hard {{ sha 'commit 4b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-5                           |
      | branch-5 | git reset --hard {{ sha 'commit 5b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
