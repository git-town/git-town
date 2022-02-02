Feature: display the main branch configuration

  Scenario: not configured
    Given my repo doesn't have a main branch configured
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
