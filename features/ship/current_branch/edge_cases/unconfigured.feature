@skipWindows
Feature: ask for missing configuration information

  Scenario: running unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town ship" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now "main"
    And my repo is now configured with no perennial branches
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
