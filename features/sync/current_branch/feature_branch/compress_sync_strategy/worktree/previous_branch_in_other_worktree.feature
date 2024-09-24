Feature: sync while the previous branch is checked out in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE   |
      | current | local    | current 1 |
      | current | local    | current 2 |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    And Git Town setting "sync-strategy" is "compress"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | current | git fetch --prune --tags      |
      |         | git checkout main             |
      | main    | git rebase origin/main        |
      |         | git checkout current          |
      | current | git merge --no-edit --ff main |
      |         | git reset --soft main         |
      |         | git commit -m "current 1"     |
      |         | git push -u origin current    |
    And the current branch is still "current"
    And the previous Git branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                |
      | current | git reset --hard {{ sha 'current 2' }} |
      |         | git push origin :current               |
    And the current branch is now "current"
    And the previous Git branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
