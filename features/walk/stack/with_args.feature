Feature: walk each branch of a stack with arguments

  @debug @this
  Scenario: iterate the full stack
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --stack echo hello"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
      | (none)   | echo hello            |
      | branch-1 | git checkout branch-2 |
      | (none)   | echo hello            |
      | branch-2 | git checkout branch-3 |
      | (none)   | echo hello            |
      | branch-3 | git checkout branch-2 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
