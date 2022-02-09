Feature: continue after successful command

  Scenario Outline:
    Given my repo has a feature branch "feature"
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
      | append new           |
      | completions fish     |
      | config               |
      | diff-parent          |
      | hack new             |
      | help                 |
      | kill feature         |
      | main_branch          |
      | new-branch-push-flag |
      | new-pull-request     |
      | offline              |
      | perennial-branches   |
      | prepend new          |
      | prune-branches       |
      | pull-branch-strategy |
      | rename-branch        |
      | repo                 |
      | ship feature -m done |
      | sync                 |
      | version              |
