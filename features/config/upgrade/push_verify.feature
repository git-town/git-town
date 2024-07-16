Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo clone
    And <LOCATION> Git Town setting "push-verify" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.push-verify" to "git-town.push-hook".
      """
    And <LOCATION> Git Town setting "push-hook" is now "true"
    And <LOCATION> Git Town setting "push-verify" now doesn't exist

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
