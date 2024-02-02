Feature: automatically upgrade outdated configuration

  @this
  Scenario Outline:
    Given <LOCATION> Git Town setting "code-hosting-origin-hostname" is "git.acme.com"
    When I run "git-town config"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.code-hosting-origin-hostname".
      I am upgrading this setting to the new format "git-town.hosting-origin-hostname".
      """
    And <LOCATION> Git Town setting "hosting-origin-hostname" is now "git.acme.com"
    And <LOCATION> Git Town setting "code-hosting-origin-hostname" now doesn't exist

    Examples:
      | LOCATION |
      # | local    |
      | global   |
