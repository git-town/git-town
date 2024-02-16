Feature: rename the current branch to a branch that is active in another worktree

  Background:
    Given the feature branches "old" and "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | main commit  |
      | old    | local, origin | old commit   |
      | other  | local         | other commit |
    And branch "other" is active in another worktree
    And the current branch is "old"
    When I run "git-town rename-branch other"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "other"
      """
    And the current branch is still "old"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is now "old"
