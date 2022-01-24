Feature: display the main branch configuration

  As a user or tool unsure about which branch is currently configured as the main branch
  I want to be able to see this information simply and directly
  So that I can use it without furter thinking or processing, and my Git Town workflows are effective.

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
