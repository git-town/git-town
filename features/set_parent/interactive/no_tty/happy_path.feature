@skipWindows
Feature: no TTY

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
    And the current branch is "branch-1"
    When I run "git-town set-parent branch-2" in a non-TTY shell

  Scenario: result
    Then Git Town runs no commands
    And this lineage exists now
      """
      main
        branch-2
          branch-1
      """

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs no commands
    And this lineage exists now
      """
      main
        branch-1
        branch-2
      """
