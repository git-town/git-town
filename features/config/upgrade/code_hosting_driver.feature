Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "code-hosting-driver" is "github"
    When I run "git-town config"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.code-hosting-driver" to "git-town.hosting-platform".
      """
    And <LOCATION> Git Town setting "hosting-platform" is now "github"
    And <LOCATION> Git Town setting "code-hosting-driver" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
