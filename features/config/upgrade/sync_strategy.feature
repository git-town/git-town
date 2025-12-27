Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.sync-strategy" is "rebase"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.sync-strategy to git-town.sync-feature-strategy.
      """
    And <LOCATION> Git setting "git-town.sync-feature-strategy" is now "rebase"
    And <LOCATION> Git setting "git-town.sync-strategy" now doesn't exist

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
