Feature: too few or many parameters

  Scenario Outline:
    When I run "git-town <CMD>"
    Then it runs no commands
    And it prints the error:
      """
      Usage:
      """

    Examples:
      | CMD                                   |
      | alias                                 |
      | alias arg1 arg2                       |
      | append                                |
      | append arg1 arg2                      |
      | completions arg1                      |
      | config arg1                           |
      | config main-branch arg1 arg2          |
      | config new-branch-push-flag arg1 arg2 |
      | hack                                  |
      | hack arg1 arg2                        |
      | kill arg1 arg2                        |
      | new-pull-request arg1                 |
      | offline arg1 arg2                     |
      | perennial-branches arg1               |
      | perennial-branches update arg1        |
      | prepend                               |
      | prune-branches arg1                   |
      | pull-branch-strategy arg1 arg2        |
      | pull-branch-strategy invalid          |
      | rename-branch                         |
      | rename-branch arg1 arg2 arg3          |
      | repo arg1                             |
      | set-parent-branch arg1                |
      | ship arg1 arg2                        |
      | sync arg1                             |
      | version arg1                          |
