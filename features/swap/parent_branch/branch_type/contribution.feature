@this
Feature: swapping a branch with its contribution parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT       | LOCATIONS     |
      | contribution | contribution | main         | local, origin |
      | branch-1     | feature      | contribution | local, origin |
    And the current branch is "branch-1"
    When I run "git-town swap"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch "contribution" is a contribution branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "branch-1"
    And the initial commits exist now
    And the initial lineage exists now
