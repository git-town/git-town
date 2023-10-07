Feature: append a branch to a branch whose tracking branch was deleted

  Background:
    Given the feature branch "shipped"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, origin | shipped commit |
      | shipped | local, origin | shipped commit |
    And origin ships the "shipped" branch
    And the current branch is "shipped"
    And an uncommitted file
    When I run "git-town append new"

  # @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | shipped | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
      |         | git checkout shipped     |
      | shipped | git merge --no-edit main |
      |         | git checkout main        |
      | main    | git branch -d shipped    |
      |         | git branch new main      |
      |         | git checkout new         |
      | new     | git stash pop            |
    And it prints:
      """
      Cannot append branch "new" to branch "old"
      because branch "old" has been deleted at the "origin" remote.
      """
    And the current branch is now "main"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no branch hierarchy exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git add -A                            |
      |        | git stash                             |
      |        | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
      | old    | git stash pop                         |
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial branches and hierarchy exist
