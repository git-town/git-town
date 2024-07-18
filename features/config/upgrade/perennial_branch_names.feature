Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo clone
    And <LOCATION> Git Town setting "perennial-branch-names" is "one two"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.perennial-branch-names" to "git-town.perennial-branches".
      """
    And <LOCATION> Git Town setting "perennial-branches" is now "one two"
    And <LOCATION> Git Town setting "perennial-branch-names" now doesn't exist

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
