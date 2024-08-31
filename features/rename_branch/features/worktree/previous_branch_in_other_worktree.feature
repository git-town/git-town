Feature: rename the current branch to a branch that is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git branch new current   |
      |         | git checkout new         |
      | new     | git branch -D current    |
    And the current branch is now "new"
    And the previous Git branch is now "new"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | new     | git branch current {{ sha 'initial commit' }} |
      |         | git checkout current                          |
      | current | git branch -D new                             |
    And the current branch is now "current"
    And the previous Git branch is now ""
