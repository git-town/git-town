Feature: automatically upgrade outdated configuration

  @this
  Scenario Outline:
    Given <LOCATION> Git Town setting "sync-delete-remote-branch" is "true"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.sync-delete-remote-branch".
      I am upgrading this setting to the new format "git-town.sync-delete-tracking-branch".
      """
    And <LOCATION> Git Town setting "sync-delete-tracking-branch" is now "true"
    And <LOCATION> Git Town setting "sync-delete-remote-branch" no longer exists

    Examples:
      | COMMAND | LOCATION |
      | config  | local    |
      | config  | global   |
