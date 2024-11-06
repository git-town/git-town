Feature: ship a feature branch in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | current | local    | current commit |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | current | git checkout main               |
      | main    | git merge --squash --ff current |
      |         | git commit -m "feature done"    |
      |         | git branch -D current           |
    And the current branch is now "main"
    And the previous Git branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch current {{ sha 'current commit' }} |
      |        | git checkout current                          |
    And the current branch is now "current"
    And the previous Git branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
