@skipWindows
Feature: ask for missing configuration information

  Scenario: run unconfigured
    Given Git Town is not configured
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
