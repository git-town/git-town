Feature: previous Git branch is in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town delete"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -D current    |
    And the current branch is now "main"
    And the previous Git branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch current {{ sha 'initial commit' }} |
      |        | git checkout current                          |
    And the current branch is now "current"
    And the previous Git branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
