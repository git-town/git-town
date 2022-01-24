Feature: offline mode

  Scenario: trying to create a new pull request in offline mode
    Given Git Town is in offline mode
    When I run "git-town new-pull-request"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
