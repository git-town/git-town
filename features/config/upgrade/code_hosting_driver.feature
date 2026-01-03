Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.code-hosting-driver" is "github"
    When I run "git-town config"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.code-hosting-driver to git-town.forge-type.
      """
    And <LOCATION> Git setting "git-town.code-hosting-driver" now doesn't exist
    And <LOCATION> Git setting "git-town.forge-type" is now "github"

    Examples:
      | LOCATION |
      | local    |
      | global   |
