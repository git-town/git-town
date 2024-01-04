Feature: configure the sync-perennial-strategy

  Scenario Outline:
    When I run "git-town config sync-perennial-strategy <VALUE>"
    Then local Git Town setting "sync-perennial-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-perennial-strategy rebase"
    Then the configuration file is now:
      """
      [sync-strategy]
        perennial-branches = "rebase"
      """
    And local Git Town setting "sync-perennial-strategy" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      [sync-strategy]
        perennial-branches = "rebase"
      """
    When I run "git-town config sync-perennial-strategy merge"
    Then the configuration file is now:
      """
      [sync-strategy]
        perennial-branches = "merge"
      """
    And local Git Town setting "sync-perennial-strategy" still doesn't exist
