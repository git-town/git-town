Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo clone
    And <LOCATION> Git Town setting "code-hosting-origin-hostname" is "git.acme.com"
    When I run "git-town config"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.code-hosting-origin-hostname" to "git-town.hosting-origin-hostname".
      """
    And <LOCATION> Git Town setting "hosting-origin-hostname" is now "git.acme.com"
    And <LOCATION> Git Town setting "code-hosting-origin-hostname" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
