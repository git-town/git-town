Feature: update the parent of a feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS     |
      | parent branch of child | up enter |

  Scenario: result
    Then this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
