Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git Town setting "main-branch-name" is "main"
    When I run "git-town hack foo"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.main-branch-name" to "git-town.main-branch".
      """
    And <LOCATION> Git Town setting "main-branch" is now "main"
    And <LOCATION> Git Town setting "main-branch-name" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
