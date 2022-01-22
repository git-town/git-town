@skipWindows
Feature: git town config setup

  To reliably update the Git Town configuration
  When Git Town is already configured
  I want to review and enter all core configuration values again.

  Scenario: everything is already configured
    Given my repo has the branches "production" and "qa"
    And the main branch is configured as "main"
    And the perennial branches are configured as "qa"
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [ENTER]                     |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production"

  Scenario: some configuration entries are missing
    Given my repo has the branches "production" and "dev"
    And I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [DOWN][ENTER]               |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "dev" and "production"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER        |
      | Please specify the main development branch | [DOWN][ENTER] |
    Then the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
