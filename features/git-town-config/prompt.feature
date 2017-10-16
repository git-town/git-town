Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town
  I want to be prompted to configure it when I use it the first time
  So that I use a properly configured tool at all times.


  Background:
    Given I haven't configured Git Town yet


  @ignore-run-error
  Scenario Outline: All Git Town commands show the configuration prompt if running unconfigured
    When I run `<COMMAND>` and enter "main" and ""
    Then Git Town prints the initial configuration prompt
    And the main branch is now configured as "main"
    And my repo is configured with no perennial branches

    Examples:
      | COMMAND                   |
      | git-town hack feature     |
      | git-town kill             |
      | git-town new-pull-request |
      | git-town prune-branches   |
      | git-town repo             |
      | git-town ship             |
      | git-town sync             |
