Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.code-hosting-origin-hostname" is "git.acme.com"
    When I run "git-town config"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.code-hosting-origin-hostname" to "git-town.hosting-origin-hostname".
      """
    And <LOCATION> Git setting "git-town.hosting-origin-hostname" is now "git.acme.com"
    And <LOCATION> Git setting "git-town.code-hosting-origin-hostname" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
