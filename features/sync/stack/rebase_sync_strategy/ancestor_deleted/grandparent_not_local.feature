Feature: the parent of a branch was deleted remotely and the grandparent exists only on origin

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
      | branch-2 | local, origin | branch-2 commit |
      | branch-3 | local, origin | branch-3 commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "branch-2" branch
    And the current branch is "branch-3"
    And I ran "git branch -D branch-1"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                               |
      | branch-3 | git fetch --prune --tags                                              |
      |          | git pull                                                              |
      |          | git -c rebase.updateRefs=false rebase --onto origin/branch-1 branch-2 |
      |          | git push --force-with-lease                                           |
      |          | git branch -D branch-2                                                |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-1
          branch-3
      """
