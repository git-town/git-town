@skipWindows
Feature: "git town config setup" with existing configuration

  To reliably update the Git Town configuration
  When Git Town is already configured
  I want to review and enter all core configuration values again.

  Background: everything is configured
    Given my repo has the feature branches "production" and "qa"
    And the main branch is configured as "main"
    And the perennial branches are configured as "qa"
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [ENTER]                     |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |

  Scenario: result
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production"
