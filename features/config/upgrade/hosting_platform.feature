Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.hosting-platform" is "github"
    When I run "git-town config"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.hosting-platform" to "git-town.forge-type".
      """
    And <LOCATION> Git setting "git-town.forge-type" is now "github"
    And <LOCATION> Git setting "git-town.hosting-platform" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
