Feature: automatically upgrade the codeberg-token to forgejo-token

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.codeberg-token" is "token"
    When I run "git-town hack foo"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.codeberg-token to git-town.forgejo-token.
      """
    And <LOCATION> Git setting "git-town.codeberg-token" now doesn't exist
    And <LOCATION> Git setting "git-town.forgejo-token" is now "token"

    Examples:
      | LOCATION |
      | local    |
      | global   |
