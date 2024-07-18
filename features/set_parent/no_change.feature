Feature: update the parent of a feature branch

  Background:
    Given a Git repo clone
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
    And it runs no commands
    And the initial lineage exists
    And the current branch is still "child"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
