Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "code-hosting-driver" is "github"
    When I run "git-town config"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.code-hosting-driver".
      I am upgrading this setting to the new format "git-town.code-hosting-platform".
      """
    And <LOCATION> Git Town setting "code-hosting-platform" is now "github"
    And <LOCATION> Git Town setting "code-hosting-driver" no longer exists

    Examples:
      | LOCATION |
      | local    |
      | global   |
