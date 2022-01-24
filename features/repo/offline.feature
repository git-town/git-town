Feature: git town-repo: offline mode

  When offline
  I want Git Town to tell me that I cannot see my repository
  So that I don't get misleading error messages.

  Scenario: trying to prune branches in offline mode
    Given Git Town is in offline mode
    When I run "git-town repo"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
