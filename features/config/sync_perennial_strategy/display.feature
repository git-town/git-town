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

  Scenario: empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-perennial-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario: set in config file
    Given the configuration file:
      """
      [sync-strategy]
        perennial-branches = "merge"
      """
    When I run "git-town config sync-perennial-strategy"
    Then it prints:
      """
      merge
      """

  Scenario: illegal setting in config file
    Given the configuration file:
      """
      [sync-strategy]
        perennial-branches = "zonk"
      """
    When I run "git-town config sync-perennial-strategy"
    Then it prints the error:
      """
      unknown sync-perennial strategy: "zonk"
      """
