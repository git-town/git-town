Feature: does not compress empty branches

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town compress"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" has no commits
      """
    And the current branch is still "feature"
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And the initial branches and lineage exist
