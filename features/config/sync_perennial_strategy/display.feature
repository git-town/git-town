Feature: display the currently configured sync_perennial_strategy

  Scenario: default
    When I run "git-town config sync-perennial-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario Outline: configured locally
    Given local Git Town setting "sync-perennial-strategy" is "<VALUE>"
    When I run "git-town config sync-perennial-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: configured globally
    Given global Git Town setting "sync-perennial-strategy" is "<VALUE>"
    When I run "git-town config sync-perennial-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
