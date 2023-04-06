Feature: offline mode

  Scenario:
    Given offline mode is enabled
    When I run "git-town new-pull-request"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
