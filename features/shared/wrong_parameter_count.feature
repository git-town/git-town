Feature: too few or many parameters

  Scenario Outline: incorrect number of arguments
    When I run "git-town <CMD>"
    Then it runs no commands
    And it prints the error:
      """
      Usage:
      """

    Examples:
      | CMD                                   |
      | aliases                               |
      | aliases arg1 arg2                     |
      | append                                |
      | append arg1 arg2                      |
      | completions arg1                      |
      | config arg1                           |
      | config main-branch arg1 arg2          |
      | config push-new-branches arg1 arg2    |
      | config offline arg1 arg2              |
      | config perennial-branches arg1        |
      | config perennial-branches update arg1 |
      | config pull-branch-strategy arg1 arg2 |
      | hack                                  |
      | hack arg1 arg2                        |
      | kill arg1 arg2                        |
      | new-pull-request arg1                 |
      | prepend                               |
      | prune-branches arg1                   |
      | rename-branch                         |
      | rename-branch arg1 arg2 arg3          |
      | repo arg1                             |
      | set-parent arg1                       |
      | ship arg1 arg2                        |
      | sync arg1                             |
      | version arg1                          |

  Scenario Outline: invalid arguments
    When I run "git-town <CMD>"
    Then it runs no commands
    And it prints the error:
      """
      Error: <ERROR>
      """

    Examples:
      | CMD                                 | ERROR                                   |
      | config pull-branch-strategy invalid | unknown pull branch strategy: "invalid" |
      | config sync-strategy invalid        | unknown sync strategy: "invalid"        |
