Feature: display the currently configured sync-strategy

  Scenario: default
    When I run "git-town config sync-strategy"
    Then it prints:
      """
      merge
      """

  Scenario Outline: local setting
    Given local setting "sync-strategy" is "<VALUE>"
    When I run "git-town config sync-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: global setting
    Given global setting "sync-strategy" is "<VALUE>"
    When I run "git-town config sync-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: global and local set to different values
    Given global setting "sync-strategy" is "merge"
    And local setting "sync-strategy" is "rebase"
    When I run "git-town config sync-strategy <FLAG>"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | FLAG     | OUTPUT |
      | --global | merge  |
      |          | rebase |
