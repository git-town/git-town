@messyoutput
Feature: select the new parent via a visual dialog

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the current branch is "branch-2"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                       | KEYS       |
      | parent branch for "branch-2" | down enter |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "branch-2" is now a child of "branch-1"
      """
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
