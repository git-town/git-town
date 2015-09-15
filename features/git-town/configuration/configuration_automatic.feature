Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town yet
  I want to get a friendly reminder with the opportunity to configure it right now when I use it the first time
  So that I use a properly configured tool at all times.


  Background:
    Given I haven't configured Git Town yet


  Scenario Outline: All Git Town commands show the configuration prompt if running unconfigured
    When I run `<COMMAND>` and enter "main" and ""
    Then I see the initial configuration prompt
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none
    And it may error

    Examples:
      | COMMAND              |
      | git extract          |
      | git hack             |
      | git kill             |
      | git new-pull-request |
      | git prune-branches   |
      | git repo             |
      | git ship             |
      | git sync             |
      | git sync-fork        |
