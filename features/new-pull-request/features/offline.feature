Feature: offline mode

  Scenario:
    Given Git Town is in offline mode
    When I run "git-town new-pull-request"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
