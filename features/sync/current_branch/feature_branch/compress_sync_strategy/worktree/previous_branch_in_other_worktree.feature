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
    And Git setting "git-town.sync-strategy" is "compress"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                    |
      | current | git fetch --prune --tags   |
      |         | git reset --soft main      |
      |         | git commit -m "current 1"  |
      |         | git push -u origin current |
    And the previous Git branch is still "previous"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                |
      | current | git reset --hard {{ sha 'current 2' }} |
      |         | git push origin :current               |
    And the previous Git branch is still "previous"
    And the initial commits exist now
    And the initial branches and lineage exist now
