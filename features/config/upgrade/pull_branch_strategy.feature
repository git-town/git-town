Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "pull-branch-strategy" is "rebase"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.pull-branch-strategy".
      I am upgrading this setting to the new format "git-town.sync-perennial-strategy".
      """
    And <LOCATION> Git Town setting "sync-perennial-strategy" is now "rebase"
    And <LOCATION> Git Town setting "pull-branch-strategy" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
