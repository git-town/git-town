Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.ship-delete-remote-branch" is "true"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.ship-delete-remote-branch to git-town.ship-delete-tracking-branch.
      """
    And <LOCATION> Git setting "git-town.ship-delete-remote-branch" now doesn't exist
    And <LOCATION> Git setting "git-town.ship-delete-tracking-branch" is now "true"

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
