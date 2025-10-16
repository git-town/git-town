@messyoutput
Feature: update the parent of a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                    | KEYS  |
      | parent branch for "child" | enter |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "child" is now a child of "parent"
      """
    And the initial lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
