Feature: ship while the previous branch is active in another worktree

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | current | local    | current commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | current | git checkout main                   |
      | main    | git merge --no-ff --edit -- current |
      |         | git branch -D current               |
    And the previous Git branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch current {{ sha 'current commit' }} |
      |        | git checkout current                          |
    And the previous Git branch is now "main"
