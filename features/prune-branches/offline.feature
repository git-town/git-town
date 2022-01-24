Feature: git town-prune-branches: offline mode

  When offline
  I want to still be able to prune branches
  So that I can work as much as possible even without internet connection.

  Scenario: trying to prune branches in offline mode
    Given Git Town is in offline mode
    When I run "git-town prune-branches"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
