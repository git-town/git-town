@skipWindows
Feature: Ask for missing configuration information

  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now "main"
    And my repo is now configured with no perennial branches
