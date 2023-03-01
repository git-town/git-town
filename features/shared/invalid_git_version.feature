Feature: require minimum Git version

  Scenario Outline: using an unsupported Git Version
    Given Git has version "2.6.2"
    When I run "git-town <COMMAND>"
    Then it prints the error:
      """
      this app requires Git 2.7.0 or higher
      """

    Examples:
      | COMMAND              |
      |                      |
      | append               |
      | config               |
      | diff-parent          |
      | hack                 |
      | help                 |
      | install aliases true |
      | kill                 |
      | main_branch          |
      | push-new-branches    |
      | new-pull-request     |
      | offline              |
      | perennial-branches   |
      | prepend              |
      | prune-branches       |
      | pull-branch-strategy |
      | rename-branch        |
      | repo                 |
      | set-parent           |
      | ship                 |
      | sync                 |
      | version              |
