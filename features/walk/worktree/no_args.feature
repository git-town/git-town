Feature: walk each branch of a stack without arguments when some branches are checked out in another worktree

  @this
  Scenario: action
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
    And the current branch is "branch-2"
    And branch "branch-1" is active in another worktree
    When I run "git-town walk --stack"
    Then Git Town runs no commands
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And the current branch is still "branch-2"
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-3 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And the current branch is now "branch-3"
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-3 | git checkout branch-2 |
    And Git Town prints:
      """
      Branch walk done.
      """
