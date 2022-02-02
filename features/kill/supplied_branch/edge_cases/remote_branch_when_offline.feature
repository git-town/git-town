Feature: cannot kill a remote branch in offline mode

  Background:
    Given Git Town is in offline mode
    And my origin has a feature branch "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | remote   | feature commit |
    And my repo knows about the remote branch
    And I am on the "main" branch
    When I run "git-town kill feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND |
    And it prints the error:
      """
      cannot delete remote branch "feature" in offline mode
      """
    And I am still on the "main" branch
    And Git Town still has no branch hierarchy information
