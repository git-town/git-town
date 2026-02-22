@skipWindows
Feature: no TTY, unknown parent

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | branch | (none) |        | local, origin |
    And the current branch is "branch"
    When I run "git-town ship" in a non-TTY shell
    Then Git Town prints the error:
      """
      cannot determine parent branch for "branch": no interactive terminal available

      To configure, run:
      git checkout branch && git-town set-parent <parent-branch>
      """
