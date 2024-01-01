Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "main-branch-name" is "main"
    When I run "git-town hack foo"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.main-branch-name".
      I am upgrading this setting to the new format "git-town.main-branch".
      """
    And <LOCATION> Git Town setting "main-branch" is now "main"
    And <LOCATION> Git Town setting "main-branch-name" no longer exists

    Examples:
      | LOCATION |
      | local    |
      | global   |
