Feature: git town-prune-branches: offline mode


  Scenario: trying to prune branches in offline mode
    Given Git Town is in offline mode
    When I run "git-town prune-branches"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
