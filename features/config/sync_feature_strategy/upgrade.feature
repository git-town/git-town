Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "sync-strategy" is "rebase"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.sync-strategy".
      I am upgrading this setting to the new format "git-town.sync-feature-strategy".
      """
    And <LOCATION> Git Town setting "sync-feature-strategy" is now "rebase"
    And <LOCATION> Git Town setting "sync-strategy" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
