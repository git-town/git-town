Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.perennial-branch-names" is "one two"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.perennial-branch-names to git-town.perennial-branches.
      """
    And <LOCATION> Git setting "git-town.perennial-branch-names" now doesn't exist
    And <LOCATION> Git setting "git-town.perennial-branches" is now "one two"

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
