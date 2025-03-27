Feature: hoisting a branch out of a stack

  Background:
    Given a Git repo with origin
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
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-4 | feature | branch-3 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-4 | local    | commit 4a |
      | branch-4 | local    | commit 4b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-5 | feature | branch-4 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-5 | local    | commit 5a |
      | branch-5 | local    | commit 5b |
    # # make branch-2 a child of main and remove commits of branch-1 from it
    # And the current branch is "branch-2"
    # When I run "git rebase --onto main branch-1"
    # # remove commits of branch-2 from branch-3
    # And I run "git checkout branch-3"
    # And I run "git rebase --onto branch-1 branch-2"
    # # remove commits of branch-2 from branch-4
    # And I run "git checkout branch-4"
    # And I run "git rebase --onto branch-3 branch-2"
    # # remove commits of branch-2 from branch-5
    # And I run "git checkout branch-5"
    # And I run "git rebase --onto branch-4 branch-2"
    # # manually update the lineage
    # And I run "git config git-town-branch.branch-2.parent main"
    # And I run "git config git-town-branch.branch-3.parent branch-1"
    And the current branch is "branch-2"
    When I run "git-town hoist"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                             |
      | branch-2 | git fetch --prune --tags            |
      |          | git rebase --onto main branch-1     |
      |          | git checkout branch-3               |
      | branch-3 | git rebase --onto branch-1 branch-2 |
      |          | git checkout branch-4               |
      | branch-4 | git rebase --onto branch-3 branch-2 |
      |          | git checkout branch-5               |
      | branch-5 | git rebase --onto branch-4 branch-2 |
      |          | git checkout branch-2               |
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      |          |          | commit 1b |
      | branch-2 | local    | commit 2a |
      |          |          | commit 2b |
      | branch-3 | local    | commit 3a |
      |          |          | commit 3b |
      | branch-4 | local    | commit 4a |
      |          |          | commit 4b |
      | branch-5 | local    | commit 5a |
      |          |          | commit 5b |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
