Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.pull-branch-strategy" is "rebase"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.pull-branch-strategy to git-town.sync-perennial-strategy.
      """
    And <LOCATION> Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And <LOCATION> Git setting "git-town.pull-branch-strategy" now doesn't exist

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
