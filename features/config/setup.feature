@skipWindows
Feature: enter Git Town configuration

  Scenario: already configured
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [ENTER]                     |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now "main"
    And the perennial branches are now "production"

  Scenario: unconfigured
    Given the branches "dev" and "production"
    And Git Town is not configured
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [DOWN][ENTER]               |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now "main"
    And the perennial branches are now "dev" and "production"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER        |
      | Please specify the main development branch | [DOWN][ENTER] |
    Then the main branch is now "main"
    And there are still no perennial branches
