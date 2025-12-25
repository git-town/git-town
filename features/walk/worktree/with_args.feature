Feature: walk each branch of a stack with arguments when some branches are checked out in another worktree

  Scenario: iterate the full stack
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
    And the current branch is "branch-2"
    And branch "branch-1" is active in another worktree
    When I run "git-town walk --stack echo hello"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | echo hello            |
      |          | git checkout branch-3 |
      | branch-3 | echo hello            |
      |          | git checkout branch-2 |
    And Git Town prints:
      """
      Branch walk done.
      """
