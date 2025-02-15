Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.forge-type" is "github"
    When I run "git-town config"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.forge-type" to "git-town.forge-type".
      """
    And <LOCATION> Git setting "git-town.forge-type" is now "github"
    And <LOCATION> Git setting "git-town.forge-type" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
