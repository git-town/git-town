Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "ship-delete-remote-branch" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.ship-delete-remote-branch".
      I am upgrading this setting to the new format "git-town.ship-delete-tracking-branch".
      """
    And <LOCATION> Git Town setting "ship-delete-tracking-branch" is now "true"
    And <LOCATION> Git Town setting "ship-delete-remote-branch" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
