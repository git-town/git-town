Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town yet
  I want to get a friendly reminder with the opportunity to configure it right now when I use it the first time
  So that I use a properly configured tool at all times.


  Background:
    Given I haven't configured Git Town yet


  Scenario Outline: All Git Town commands show the configuration prompt if running unconfigured
    When I run `<COMMAND>` and enter "^C"
    Then I see the initial configuration prompt

    Examples:
    | COMMAND              |
    | git new-pull-request |
    | git extract          |
    | git hack             |
    | git kill             |
    | git prune-branches   |
    | git repo             |
    | git ship             |
    | git sync             |
    | git sync-fork        |


  Scenario Outline: Starting the wizard by answering "y" to the configuration prompt's question whether to start it
    When I run `<COMMAND>` and enter "y" and "^C"
    Then I see the first line of the configuration wizard

    Examples:
    | COMMAND              |
    | git new-pull-request |
    | git extract          |
    | git hack             |
    | git kill             |
    | git prune-branches   |
    | git repo             |
    | git ship             |
    | git sync             |
    | git sync-fork        |


  Scenario Outline: Starting the wizard by choosing the default answer to the configuration prompt's question whether to start it
    When I run `<COMMAND>` and enter "" and "^C"
    Then I see the first line of the configuration wizard

    Examples:
    | COMMAND              |
    | git new-pull-request |
    | git extract          |
    | git hack             |
    | git kill             |
    | git prune-branches   |
    | git repo             |
    | git ship             |
    | git sync             |
    | git sync-fork        |


  Scenario Outline: Not proceeding to the wizard by answering "n" to the configuration prompt's question whether to start it
    When I run `<COMMAND>` and enter "n" and "^C"
    Then I don't see the first line of the configuration wizard

    Examples:
    | COMMAND              |
    | git new-pull-request |
    | git extract          |
    | git hack             |
    | git kill             |
    | git prune-branches   |
    | git repo             |
    | git ship             |
    | git sync             |
    | git sync-fork        |
