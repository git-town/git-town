Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "ship-delete-remote-branch" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.ship-delete-remote-branch" to "git-town.ship-delete-tracking-branch".
      """
    And <LOCATION> Git Town setting "ship-delete-tracking-branch" is now "true"
    And <LOCATION> Git Town setting "ship-delete-remote-branch" now doesn't exist

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
