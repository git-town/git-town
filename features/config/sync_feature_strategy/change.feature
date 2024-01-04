Feature: configure the sync-feature-strategy

  Scenario Outline:
    When I run "git-town config sync-feature-strategy <VALUE>"
    Then local Git Town setting "sync-feature-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-feature-strategy rebase"
    Then the configuration file is now:
      """
      [sync-strategy]
        feature-branches = "rebase"
      """
    And local Git Town setting "sync-feature-strategy" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      [sync-strategy]
        feature-branches = "rebase"
      """
    When I run "git-town config sync-feature-strategy merge"
    Then the configuration file is now:
      """
      [sync-strategy]
        feature-branches = "merge"
      """
    And local Git Town setting "sync-feature-strategy" still doesn't exist
