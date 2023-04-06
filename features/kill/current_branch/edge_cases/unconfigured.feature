@skipWindows
Feature: ask for missing configuration

  Scenario:
    Given Git Town is not configured
    When I run "git-town kill" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the error:
      """
      you can only kill feature branches
      """
    And the main branch is now "main"
