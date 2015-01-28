Feature: Initial configuration

  Scenario Outline: Running Git Town commands while Git Town is unconfigured
    Given I haven't configured Git Town yet
    When I run `<COMMAND>` while allowing errors
    Then I see
      """
      Git Town hasn't been configured for this repository.
      Please run 'git town config --setup'.
      Would you like to do that now? Y/n
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
