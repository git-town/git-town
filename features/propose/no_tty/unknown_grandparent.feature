@skipWindows
Feature: no TTY, unknown parent

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | branch-1 | (none) |        | local, origin |
      | branch-2 | (none) |        | local, origin |
    And the current branch is "branch-2"
    And I ran "git-town set-parent branch-1"
    When I run "git-town propose" in a non-TTY shell
    Then Git Town prints the error:
      """
      cannot determine parent branch for "branch-1": no interactive terminal available

      To configure, run:
      git checkout branch-1 && git-town set-parent <parent-branch>
      """
