Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town
  I want to be prompted to configure it when I use it the first time
  So that I use a properly configured tool at all times.


  Scenario Outline: All Git Town commands show the configuration prompt if running unconfigured
    Given I haven't configured Git Town yet
    When I run "<COMMAND>" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
      | Please specify perennial branches          | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now configured as "main"
    And my repo is now configured with no perennial branches

    Examples:
      | COMMAND                   |
      | git-town hack feature     |
      | git-town kill             |
      | git-town new-pull-request |
      | git-town prune-branches   |
      | git-town repo             |
      | git-town ship             |
      | git-town sync             |
