Feature: display the main branch configuration

  Scenario: not configured
    Given the main branch is not set
    When I run "git-town main-branch"
    Then it prints:
      """
      [none]
      """

  Scenario: configured
    Given the main branch is "main"
    When I run "git-town main-branch"
    Then it prints:
      """
      main
      """
