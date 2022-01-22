@skipWindows
Feature: "git town config setup" without existing configuration and branches

  To reliably configure Git Town
  I want to be asked to enter all core configuration values.

  Background:
    Given my repo has the feature branches "production" and "dev"
    And I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [DOWN][ENTER]               |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |

  Scenario: result
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "dev" and "production"
