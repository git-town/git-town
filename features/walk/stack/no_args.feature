Feature: walk each branch of a stack without arguments

  Scenario: action
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
      | branch-A | feature | main     | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --stack"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And the current branch is now "branch-1"
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-1 | git checkout branch-2 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And the current branch is now "branch-2"
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
