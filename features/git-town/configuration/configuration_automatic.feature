Feature: Automatically running the configuration wizard if Git Town is unconfigured

  As a user having forgotten to configure Git Town yet
  I want to get a friendly reminder with the opportunity to configure it right now when I use it the first time
  So that I use a properly configured tool at all times.


  Scenario Outline: Running Git Town commands before Git Town is unconfigured
    Given I haven't configured Git Town yet
    When I run `<COMMAND>`
    Then I see
      """
      Git Town hasn't been configured for this repository.
      Please run 'git town config --setup'.
      Would you like to do that now? y/n
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
