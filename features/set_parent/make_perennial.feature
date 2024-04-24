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
