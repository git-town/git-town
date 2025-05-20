Feature: execute a script in all local branches

  Scenario: iterate the full stack
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
      | branch-A | feature | main     | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all echo hello"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
      | branch-1 | echo hello            |
      |          | git checkout branch-2 |
      | branch-2 | echo hello            |
      |          | git checkout branch-3 |
      | branch-3 | echo hello            |
      |          | git checkout branch-A |
      | branch-A | echo hello            |
      |          | git checkout branch-2 |
    And Git Town prints:
      """
      Branch walk done.
      """
