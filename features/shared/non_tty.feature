Feature: non-TTY usage

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS     |
      | feature | (none) |        | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE |
      | feature | local, origin | commit  |
    And the current branch is "feature"
    When I run "git-town append new" in a non-TTY shell

  @this
  Scenario: result
    Then Git Town prints the error:
      """
      could not open a new TTY: open /dev/tty: no such device or address
      """
