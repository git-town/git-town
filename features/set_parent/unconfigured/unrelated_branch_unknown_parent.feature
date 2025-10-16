Feature: an unrelated branch has an unknown parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | branch-1  | feature | main   | local     |
      | branch-2  | feature | main   | local     |
      | unrelated | (none)  |        | local     |
    And the current branch is "branch-2"
    When I run "git-town set-parent branch-1"

  Scenario: result
    And Git Town runs no commands
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
