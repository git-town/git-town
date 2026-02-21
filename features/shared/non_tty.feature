Feature: non-TTY usage

  Scenario Outline:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS     |
      | feature | (none) |        | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE |
      | feature | local, origin | commit  |
    And the current branch is "feature"
    When I run "git-town append new" in a non-TTY shell
    Then Git Town prints the error:
      """
      no interactive terminal available
      """

    @this
    Examples:
      | COMMAND    |
      | append new |
