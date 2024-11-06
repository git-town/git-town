Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git Town setting "sync-before-ship" is "true"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Deleting obsolete setting "git-town.sync-before-ship"
      """
    And <LOCATION> Git Town setting "sync-before-ship" now doesn't exist

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
