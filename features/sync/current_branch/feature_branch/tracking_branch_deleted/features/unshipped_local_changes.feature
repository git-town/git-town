Feature: sync a branch with unshipped local changesn whose tracking branch was deleted

  Background:
    Given the feature branch "shipped"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          |
      | shipped | local, origin | shipped commit   |
      |         | local         | unshipped commit |
    And origin deletes the "shipped" branch
    And the current branch is "shipped"
    And an uncommitted file
    And inspect the repo
    When I run "git-town sync"

  @debug
  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | shipped | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
      |         | git checkout shipped     |
      | shipped | git merge main           |
      |         | git diff main            |
      |         | git stash pop            |
    And it prints:
      """
      Branch "shipped" was shipped on the remote but contains unshipped changes on your machine.
      """
    And the current branch is now "shipped"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
      | local         | shipped  |
    And this branch hierarchy exists now
      | shipped | main |

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
