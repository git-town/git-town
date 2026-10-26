@skipWindows
Feature: walk each branch of a stack in dry-run mode

  Scenario: iterate the full stack
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
      | branch-A | feature | main     | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all --dry-run touch file"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                       |
      | branch-2 | git checkout --quiet branch-1 |
      | branch-1 | touch file                    |
      |          | git checkout --quiet branch-2 |
      | branch-2 | touch file                    |
      |          | git checkout --quiet branch-3 |
      | branch-3 | touch file                    |
      |          | git checkout --quiet branch-A |
      | branch-A | touch file                    |
      |          | git checkout --quiet branch-2 |
    And Git Town prints:
      """
      [branch-2] git checkout --quiet branch-1
      (dry run)

      [branch-1] touch file
      (dry run)

      [branch-1] git checkout --quiet branch-2
      (dry run)

      [branch-2] touch file
      (dry run)

      [branch-2] git checkout --quiet branch-3
      (dry run)

      [branch-3] touch file
      (dry run)

      [branch-3] git checkout --quiet branch-A
      (dry run)

      [branch-A] touch file
      (dry run)

      [branch-A] git checkout --quiet branch-2
      (dry run)
      """
    And no uncommitted files exist now
