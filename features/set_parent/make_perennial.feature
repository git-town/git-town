Feature: make a feature branch perennial

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS       |
      | parent branch of child | down enter |

  Scenario: result
    Then it prints:
      """
      Selected parent branch for "child": <none> (perennial branch)
      """
    And it runs no commands
    And the perennial branches are now "child"
    And this lineage exists now
      | BRANCH | PARENT |
      | parent | main   |
    And the current branch is still "child"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "child"
    And the initial commits exist now
    And the initial branches and lineage exist now
