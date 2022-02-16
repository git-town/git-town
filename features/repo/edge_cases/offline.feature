Feature: offline mode

  Scenario: try to prune branches in offline mode
    Given offline mode is enabled
    When I run "git-town repo"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
