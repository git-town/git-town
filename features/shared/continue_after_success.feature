Feature: continue after successful command

  Scenario Outline:
    Given a feature branch "feature"
    And I run "git-town <COMMAND>"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """

    Examples:
      | COMMAND              |
      |                      |
      | alias true           |
      | append new-feature   |
      | completions fish     |
      | config               |
      | diff-parent          |
      | hack new-feature     |
      | help                 |
      | kill feature         |
      | main_branch          |
      | new-branch-push-flag |
      | new-pull-request     |
      | offline              |
      | perennial-branches   |
      | prepend new-feature  |
      | prune-branches       |
      | pull-branch-strategy |
      | rename-branch        |
      | repo                 |
      | ship feature -m done |
      | sync                 |
      | version              |
