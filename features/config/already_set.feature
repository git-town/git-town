Feature: listing the configuration

  As a user running the Git Town configuration wizard,
  I want to see the existing configuration values
  So that I can change it more effectively

  @skipWindows
  Scenario: everything is configured
    Given my repo has the feature branches "production" and "qa"
    And the main branch is configured as "main"
    And the perennial branches are configured as "qa"
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [ENTER]                     |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production"
