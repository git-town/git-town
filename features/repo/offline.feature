Feature: git town-repo: offline mode

  Scenario: trying to prune branches in offline mode
    Given Git Town is in offline mode
    When I run "git-town repo"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
