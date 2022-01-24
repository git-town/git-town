Feature: show help screen for commands

  Scenario Outline:
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
