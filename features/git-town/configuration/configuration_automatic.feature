Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town yet
  I want to get a friendly reminder with the opportunity to configure it right now when I use it the first time
  So that I use a properly configured tool at all times.


  Background:
    Given I haven't configured Git Town yet


  Scenario Outline: Proceeding to configuration upon initial config prompt
    When I run `<COMMAND>` and enter "y" and "^C"
    Then I see "Git Town hasn't been configured for this repository."
    And I see "Please specify the main dev branch"

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git pr             |
    | git prune-branches |
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
      Would you like to do that now? [y/n]
      """

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git pr             |
    | git prune-branches |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |
