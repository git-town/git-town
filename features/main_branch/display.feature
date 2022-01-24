Feature: display the main branch configuration

  Scenario: main branch not yet configured
    Given my repo doesn't have a main branch configured
    When I run "git-town main-branch"
    Then it prints:
      """
      [none]
      """

  Scenario: main branch is configured
    Given the main branch is configured as "main"
    When I run "git-town main-branch"
    Then it prints:
      """
      main
      """
