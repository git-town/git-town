Feature: too few or many parameters

  Scenario Outline:
    Given I am outside a Git repo
    When I run "git-town <CMD>"
    Then it runs no commands
    And it prints the error:
      """
      Error: <ERROR>
      """

    Examples:
      | CMD                          | ERROR                                              |
      | append                       | accepts 1 arg(s), received 0                       |
      | append arg1 arg2             | accepts 1 arg(s), received 2                       |
      | completions arg1             | unknown completion type: "arg1"                    |
      | config arg1                  | unknown command "arg1" for "git-town config"       |
      | config setup arg1            | unknown command "arg1" for "git-town config setup" |
      | kill arg1 arg2               | accepts at most 1 arg(s), received 2               |
      | offline arg1 arg2            | accepts at most 1 arg(s), received 2               |
      | propose arg1                 | unknown command "arg1" for "git-town propose"      |
      | prepend                      | accepts 1 arg(s), received 0                       |
      | rename-branch                | accepts between 1 and 2 arg(s), received 0         |
      | rename-branch arg1 arg2 arg3 | accepts between 1 and 2 arg(s), received 3         |
      | repo arg1 arg2               | accepts at most 1 arg(s), received 2               |
      | set-parent arg1              | unknown command "arg1" for "git-town set-parent"   |
      | ship arg1 arg2               | accepts at most 1 arg(s), received 2               |
      | sync arg1                    | unknown command "arg1" for "git-town sync"         |
      | --version arg1               | unknown command "arg1" for "git-town"              |
