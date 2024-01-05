Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "sync-before-ship" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.sync-before-ship".
      I am upgrading this setting to the new format "git-town.sync-before-ship".
      """
    And <LOCATION> Git Town setting "sync-before-ship" is now "true"
    And <LOCATION> Git Town setting "sync-before-ship" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
