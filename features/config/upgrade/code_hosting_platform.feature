Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "code-hosting-platform" is "github"
    When I run "git-town config"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.code-hosting-platform".
      I am upgrading this setting to the new format "git-town.hosting-platform".
      """
    And <LOCATION> Git Town setting "hosting-platform" is now "github"
    And <LOCATION> Git Town setting "code-hosting-platform" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
