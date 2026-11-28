Feature: compress rejects conflicting commit message options

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
    And the current branch is "feature"
    When I run "git-town compress -m compressed --keep-commit-messages"

  Scenario: using --message and --keep-commit-messages together
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      the --message and --keep-commit-messages flags cannot be used together
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
