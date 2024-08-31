Feature: previous branch is checked out in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | current | git fetch --prune --tags      |
      |         | git checkout main             |
      | main    | git rebase origin/main        |
      |         | git checkout current          |
      | current | git merge --no-edit --ff main |
      |         | git push -u origin current    |
      |         | git checkout -b new main      |
    And the current branch is now "new"
    And the previous Git branch is now "new"
    And this lineage exists now
      | BRANCH   | PARENT |
      | current  | new    |
      | new      | main   |
      | previous | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | new     | git push origin :current |
      |         | git checkout current     |
      | current | git branch -D new        |
    And the current branch is now "current"
    And the previous Git branch is now ""
    And the initial commits exist
    And the initial lineage exists
