Feature: display the currently configured pull_branch_strategy

  Scenario: default
    When I run "git-town pull-branch-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario Outline:
    Given the "pull-branch-strategy" setting is "<VALUE>"
    When I run "git-town pull-branch-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
