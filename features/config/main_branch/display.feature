Feature: display the main branch configuration

  Scenario: not configured
    Given local Git Town setting "main-branch" is ""
    When I run "git-town config main-branch"
    Then it prints:
      """
      (not set)
      """

  Scenario: configured locally
    Given local Git Town setting "main-branch" is "main"
    When I run "git-town config main-branch"
    Then it prints:
      """
      main
      """

  Scenario: configured globally
    Given global Git Town setting "main-branch" is "main"
    When I run "git-town config main-branch"
    Then it prints:
      """
      main
      """
