Feature: setting the parent to a branch whose parent is unknown

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE   | PARENT | LOCATIONS |
      | branch-1  | (none) |        | local     |
      | branch-2  | (none) |        | local     |
      | unrelated | (none) |        | local     |
    And the current branch is "branch-2"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                    | KEYS       |
      | parent branch of branch-2 | down enter |

  Scenario: result
    Then it prints:
      """
      Selected parent branch for "branch-2": branch-1
      """
    And it runs no commands
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-2 | branch-1 |
    And the current branch is still "branch-2"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "branch-2"
    And the initial commits exist
    And the initial branches and lineage exist
