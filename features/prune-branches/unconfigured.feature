Feature: Ask for missing configuration

  As a user having forgotten to configure Git Town
  I want to be prompted to configure it when I use it the first time
  So that I use a properly configured tool at all times.

  @skipWindows
  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town prune-branches" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
