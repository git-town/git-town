@skipWindows
Feature: prompt for parent branch when unknown

  Scenario Outline:
    Given my repo has a branch "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town <COMMAND>" and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |

    Examples:
      | COMMAND           |
      | append feature-2  |
      | diff-parent       |
      | kill feature-1    |
      | prepend feature-2 |
      | sync              |
