Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.push-verify" is "true"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.push-verify" to "git-town.push-hook".
      """
    And <LOCATION> Git setting "git-town.push-hook" is now "true"
    And <LOCATION> Git setting "git-town.push-verify" now doesn't exist

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
