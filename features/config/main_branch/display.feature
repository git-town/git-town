Feature: display the main branch configuration

  Scenario: not configured
    Given local setting "main-branch-name" is ""
    When I run "git-town config main-branch"
    Then it prints:
      """
      (not set)
      """

  Scenario: configured locally
    Given local setting "main-branch-name" is "main"
    When I run "git-town config main-branch"
    Then it prints:
      """
      main
      """

  Scenario: configured globally
    Given global setting "main-branch-name" is "main"
    When I run "git-town config main-branch"
    Then it prints:
      """
      main
      """
