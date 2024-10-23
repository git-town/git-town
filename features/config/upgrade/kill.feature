Feature: automatically upgrade kill alias

  Scenario:
    Given a Git repo with origin
    And global Git setting "alias.kill" is "town kill"
    When I run "git town hack feat/upgrade-alias"
    Then it prints:
      """
      Upgrading deprecated global setting "alias.kill" to "alias.delete".
      """
    And global Git setting "alias.delete" is now "town delete"
    And global Git setting "alias.kill" now doesn't exist
