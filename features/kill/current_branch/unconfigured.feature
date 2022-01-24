Feature: Ask for missing configuration

  @skipWindows
  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town kill" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
    And it prints the error:
      """
      you can only kill feature branches
      """
