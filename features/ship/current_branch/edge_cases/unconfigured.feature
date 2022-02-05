@skipWindows
Feature: ask for missing configuration information

  Scenario: unconfigured
    Given Git Town is not configured
    When I run "git-town ship" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And my repo now has no perennial branches
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
