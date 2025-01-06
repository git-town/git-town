Feature: rename the current branch to a branch that is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | old   | feature | main   | local, origin |
      | other | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | main commit  |
      | old    | local, origin | old commit   |
      | other  | local         | other commit |
    And branch "other" is active in another worktree
    And the current branch is "old"
    When I run "git-town rename other"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is already a branch "other"
      """
    And the current branch is still "old"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is now "old"
