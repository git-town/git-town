Feature: automatically upgrade kill alias

  Scenario: alias set up by Git Town
    Given a Git repo with origin
    And global Git setting "alias.kill" is "town kill"
    When I run "git-town hack foo"
    Then Git Town prints:
      """
      Upgrading deprecated global setting alias.kill to alias.delete.
      """
    And global Git setting "alias.delete" is now "town delete"
    And global Git setting "alias.kill" now doesn't exist
