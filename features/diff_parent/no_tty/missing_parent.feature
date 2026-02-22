@skipWindows
Feature: no TTY, missing parent branch

  Scenario: feature branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent" in a non-TTY shell
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot determine parent branch for "feature": no interactive terminal available

      To configure, run:
      git checkout feature && git-town set-parent <parent-branch>
      """
