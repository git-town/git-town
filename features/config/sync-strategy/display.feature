Feature: display the currently configured sync-strategy

  Scenario: default
    When I run "git-town config sync-strategy"
    Then it prints:
      """
      merge
      """

  Scenario Outline:
    Given setting "sync-strategy" is "<VALUE>"
    When I run "git-town config sync-strategy"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
