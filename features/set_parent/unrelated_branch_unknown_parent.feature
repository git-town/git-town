@messyoutput
Feature: an unrelated branch has an unknown parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | branch-1  | feature | main   | local     |
      | branch-2  | feature | main   | local     |
      | unrelated | (none)  |        | local     |
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
      | branch-1 | main     |
      | branch-2 | branch-1 |
    And the current branch is still "branch-2"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial branches and lineage exist now
