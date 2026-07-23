Feature: automatically upgrade the bitbucket-app-password to bitbucket-api-token

  Scenario Outline:
    Given a Git repo with origin
    And <LOCATION> Git setting "git-town.bitbucket-app-password" is "token"
    When I run "git-town hack foo"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.bitbucket-app-password to git-town.bitbucket-api-token.
      """
    And <LOCATION> Git setting "git-town.bitbucket-api-token" is now "token"
    And <LOCATION> Git setting "git-town.bitbucket-app-password" now doesn't exist

    Examples:
      | LOCATION |
      | local    |
      | global   |
