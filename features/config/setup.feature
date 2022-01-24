@skipWindows
Feature: Entering Git Town configuration

  Scenario: already configured
    Given my repo has the branches "production" and "qa"
    And the main branch is configured as "main"
    And the perennial branches are configured as "qa"
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [ENTER]                     |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production"

  @debug
  Scenario: unconfigured
    Given I haven't configured Git Town yet
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
