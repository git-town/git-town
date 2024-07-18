Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo clone
    And <LOCATION> Git Town setting "sync-strategy" is "rebase"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.sync-strategy" to "git-town.sync-feature-strategy".
      """
    And <LOCATION> Git Town setting "sync-feature-strategy" is now "rebase"
    And <LOCATION> Git Town setting "sync-strategy" now doesn't exist

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
