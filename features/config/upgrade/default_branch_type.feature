Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.default-branch-type" is "observed"
    When I run "git-town hack foo"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.default-branch-type" to "git-town.unknown-branch-type".
      """
    And <LOCATION> Git setting "git-town.unknown-branch-type" is now "observed"
    And <LOCATION> Git setting "git-town.default-branch-type" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
