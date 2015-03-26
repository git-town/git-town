Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town yet
  I want to get a friendly reminder with the opportunity to configure it right now when I use it the first time
  So that I use a properly configured tool at all times.


  Background:
    Given I haven't configured Git Town yet


  Scenario Outline: Seeing a configuration prompt when running a Git Town command
    When I run `<COMMAND>` and enter "^C"
    Then I see "Git Town hasn't been configured for this repository."

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git prune-branches |
    | git pull-request   |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |


  Scenario Outline: Explicitly proceeding to configuration wizard upon seeing config prompt
    When I run `<COMMAND>` and enter "y" and "^C"
    Then I see the first question of the configuration wizard

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git prune-branches |
    | git pull-request   |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |


  Scenario Outline: Implicitly proceeding to configuration wizard upon seeing config prompt
    When I run `<COMMAND>` and enter "" and "^C"
    Then I see the first question of the configuration wizard

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git prune-branches |
    | git pull-request   |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |


  Scenario Outline: Not proceeding to configuration upon initial config prompt
    When I run `<COMMAND>` and enter "n" and "^C"
    Then I see
      """
      Git Town hasn't been configured for this repository.
      Please run 'git town config --setup'.
      Would you like to do that now? [Y/n]
      """

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git prune-branches |
    | git pull-request   |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |
