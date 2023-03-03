Feature: too few or many parameters

  Scenario Outline:
    When I run "git-town <CMD>"
    Then it runs no commands
    And it prints the error:
      """
      Error: <ERROR>
      """

    Examples:
      | CMD                                   | ERROR                                                                  |
      | aliases                               | accepts 1 arg(s), received 0                                           |
      | aliases arg1 arg2                     | accepts 1 arg(s), received 2                                           |
      | append                                | accepts 1 arg(s), received 0                                           |
      | append arg1 arg2                      | accepts 1 arg(s), received 2                                           |
      | completions arg1                      | unknown argument: "arg1"                                               |
      | config arg1                           | unknown command "arg1" for "git-town config"                           |
      | config main-branch arg1 arg2          | accepts at most 1 arg(s), received 2                                   |
      | config push-new-branches arg1 arg2    | accepts at most 1 arg(s), received 2                                   |
      | config offline arg1 arg2              | accepts at most 1 arg(s), received 2                                   |
      | config perennial-branches arg1        | unknown command "arg1" for "git-town config perennial-branches"        |
      | config perennial-branches update arg1 | unknown command "arg1" for "git-town config perennial-branches update" |
      | config pull-branch-strategy arg1 arg2 | accepts at most 1 arg(s), received 2                                   |
      | config pull-branch-strategy invalid   | invalid argument "invalid" for "git-town config pull-branch-strategy"  |
      | config sync-strategy invalid          | invalid value: "invalid"                                               |
      | hack                                  | accepts 1 arg(s), received 0                                           |
      | hack arg1 arg2                        | accepts 1 arg(s), received 2                                           |
      | kill arg1 arg2                        | accepts at most 1 arg(s), received 2                                   |
      | new-pull-request arg1                 | unknown command "arg1" for "git-town new-pull-request"                 |
      | prepend                               | accepts 1 arg(s), received 0                                           |
      | prune-branches arg1                   | unknown command "arg1" for "git-town prune-branches"                   |
      | rename-branch                         | accepts between 1 and 2 arg(s), received 0                             |
      | rename-branch arg1 arg2 arg3          | accepts between 1 and 2 arg(s), received 3                             |
      | repo arg1                             | unknown command "arg1" for "git-town repo"                             |
      | set-parent arg1                       | unknown command "arg1" for "git-town set-parent"                       |
      | ship arg1 arg2                        | accepts at most 1 arg(s), received 2                                   |
      | sync arg1                             | unknown command "arg1" for "git-town sync"                             |
      | version arg1                          | unknown command "arg1" for "git-town version"                          |
