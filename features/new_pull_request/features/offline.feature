Feature: offline mode

  Scenario:
    Given offline mode is enabled
    When I run "git-town propose"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
