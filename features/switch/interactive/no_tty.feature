@skipWindows
Feature: switch branches without TTY

  Scenario: switching to another branch
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    And the current branch is "alpha"
    When I run "git-town switch" in a non-TTY shell
    Then Git Town prints the error:
      """
      no interactive terminal available
      """
