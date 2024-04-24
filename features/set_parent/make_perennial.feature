Feature: make a feature branch perennial

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS       |
      | parent branch of child | down enter |

  Scenario: result
    Then it runs no commands
    And the perennial branches are now "child"
    And this lineage exists now
      | BRANCH | PARENT |
      | parent | main   |
    And the current branch is still "child"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
