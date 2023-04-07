Feature: display the currently configured pull_branch_strategy

  Scenario: default
    When I run "git-town config pull-branch-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario Outline: configured locally
    Given local setting "pull-branch-strategy" is "<VALUE>"
    When I run "git-town config pull-branch-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |

  Scenario Outline: configured globally
    Given global setting "pull-branch-strategy" is "<VALUE>"
    When I run "git-town config pull-branch-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
