Feature: automatically upgrade rename-branch alias

  Scenario:
    Given a Git repo with origin
    And global Git setting "alias.rename-branch" is "town rename-branch"
    When I run "git town hack feat/upgrade-alias"
    Then it prints:
      """
      Upgrading deprecated global setting "alias.rename-branch" to "alias.rename".
      """
    And global Git setting "alias.rename" is now "town rename"
    And global Git setting "alias.rename-branch" now doesn't exist
