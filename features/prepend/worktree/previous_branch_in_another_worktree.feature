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
    Then Git Town runs the commands
      | BRANCH  | COMMAND                    |
      | current | git fetch --prune --tags   |
      |         | git push -u origin current |
      |         | git checkout -b new main   |
    And the previous Git branch is now "new"
    And this lineage exists now
      """
      main
        new
          current
        previous
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | new     | git checkout current     |
      | current | git branch -D new        |
      |         | git push origin :current |
    And the initial lineage exists now
    And there is now no previous Git branch
    And the initial commits exist now
