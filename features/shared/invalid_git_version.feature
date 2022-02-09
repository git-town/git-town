Feature: require minimum Git version

  Scenario Outline: using an unsupported Git Version
    Given Git has version "2.6.2"
    When I run "git-town <COMMAND>"
    Then it prints the error:
      """
      Git Town requires Git 2.7.0 or higher
      """

    Examples:
      | COMMAND              |
      |                      |
      | alias true           |
      | append               |
      | config               |
      | diff-parent          |
      | hack                 |
      | help                 |
      | kill                 |
      | main_branch          |
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
