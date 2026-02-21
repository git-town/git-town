Feature: non-TTY usage

  Scenario Outline:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | branch-1 | (none) |        | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And the current branch is "<BRANCH>"
    When I run "git-town append new" in a non-TTY shell
    Then Git Town prints the error:
      """
      no interactive terminal available
      """

    @this
    Examples:
      | BRANCH   | COMMAND    |
      | branch-1 | append new |
      | branch-1 | branch     |
