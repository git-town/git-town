Feature: help for commands

  Scenario Outline: known commands
    When I run "git-town help <COMMAND>"
    Then it prints:
      """
      Usage:
        git-town <COMMAND>
      """

    Examples:
      | COMMAND                     |
      | aliases                     |
      | append                      |
      | completions                 |
      | config                      |
      | config main-branch          |
      | config push-new-branches    |
      | config offline              |
      | config perennial-branches   |
      | config pull-branch-strategy |
      | config sync-strategy        |
      | diff-parent                 |
      | hack                        |
      | help                        |
      | kill                        |
      | new-pull-request            |
      | prepend                     |
      | prune-branches              |
      | rename-branch               |
      | repo                        |
      | set-parent                  |
      | ship                        |
      | sync                        |
      | version                     |

  Scenario Outline: outside a Git repository
    Given I am outside a Git repo
    When I run "<COMMAND>"
    Then it prints:
      """
      Usage:
        git-town [command]
      """
    And it does not print "fatal: Not a Git repository"

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
