Feature: update the parent of a feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"

  Scenario: select the default branch (current parent)
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS  |
      | parent branch of child | enter |
    And the initial lineage exists

  Scenario: select another branch
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS     |
      | parent branch of child | up enter |
    Then this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS       |
      | parent branch of child | down enter |
    Then the perennial branches are now "child"
    And this lineage exists now
      | BRANCH | PARENT |
      | parent | main   |
