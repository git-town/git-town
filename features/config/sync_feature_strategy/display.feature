Feature: display the currently configured sync-feature-strategy

  Scenario: default
    When I run "git-town config sync-feature-strategy"
    Then it prints:
      """
      merge
      """

  Scenario Outline: configured in local Git metadata
    Given local Git Town setting "sync-feature-strategy" is "<VALUE>"
    When I run "git-town config sync-feature-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: configured in global Git metadata
    Given global Git Town setting "sync-feature-strategy" is "<VALUE>"
    When I run "git-town config sync-feature-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: global and local set to different values
    Given global Git Town setting "sync-feature-strategy" is "merge"
    And local Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town config sync-feature-strategy <FLAG>"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | FLAG     | OUTPUT |
      | --global | merge  |
      |          | rebase |

  Scenario: empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-feature-strategy"
    Then it prints:
      """
      merge
      """

  Scenario: set in config file
    Given the configuration file:
      """
      [sync-strategy]
        feature-branches = "rebase"
      """
    When I run "git-town config sync-feature-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario: illegal setting in config file
    Given the configuration file:
      """
      [sync-strategy]
        feature-branches = "zonk"
      """
    When I run "git-town config sync-feature-strategy"
    Then it prints the error:
      """
      unknown sync-feature strategy: "zonk"
      """
