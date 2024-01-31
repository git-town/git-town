Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "perennial-branch-names" is "one two"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.perennial-branch-names".
      I am upgrading this setting to the new format "git-town.perennial-branches".
      """
    And <LOCATION> Git Town setting "perennial-branches" is now "one two"
    And <LOCATION> Git Town setting "perennial-branch-names" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
