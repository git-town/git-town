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
      | DIALOG                 | KEYS  |
      | parent branch of child | enter |

  Scenario: result
    Then it prints:
      """
      Selected parent branch for "child": parent
      """
    And Git Town runs no commands
    And the initial lineage exists now
    And the current branch is still "child"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "child"
    And the initial commits exist now
    And the initial branches and lineage exist now
