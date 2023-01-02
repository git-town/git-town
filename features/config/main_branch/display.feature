Feature: display the main branch configuration

  Scenario: not configured
    Given the main branch is not set
    When I run "git-town config main-branch"
    Then it prints:
      """
      (not set)
      """

  Scenario: configured
    Given the main branch is "main"
    When I run "git-town config main-branch"
    Then it prints:
      """
      main
      """
