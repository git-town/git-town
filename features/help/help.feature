Feature: show help screen for commands

  Scenario Outline: known commands
    When I run "git-town help <COMMAND>"
    Then it prints:
      """
      Usage:
        git-town <COMMAND>
      """

    Examples:
      | COMMAND              |
      | alias                |
      | append               |
      | completions          |
      | config               |
      | diff-parent          |
      | hack                 |
      | help                 |
      | kill                 |
      | main-branch          |
      | new-branch-push-flag |
      | new-pull-request     |
      | offline              |
      | perennial-branches   |
      | prepend              |
      | prune-branches       |
      | pull-branch-strategy |
      | rename-branch        |
      | repo                 |
      | set-parent-branch    |
      | ship                 |
      | sync                 |
      | version              |

  Scenario Outline: Running outside of a Git repository
    Given my workspace is currently not a Git repo
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
