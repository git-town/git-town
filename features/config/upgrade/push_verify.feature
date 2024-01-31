Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "push-verify" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.push-verify".
      I am upgrading this setting to the new format "git-town.push-hook".
      """
    And <LOCATION> Git Town setting "push-hook" is now "true"
    And <LOCATION> Git Town setting "push-verify" no longer exists

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
