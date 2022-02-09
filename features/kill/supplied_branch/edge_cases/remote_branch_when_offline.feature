Feature: does not kill a remote branch in offline mode

  Background:
    Given a remote feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | origin   | feature commit |
    And I fetch Git updates
    And offline mode is enabled
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
    And Git Town is still aware of no branch hierarchy
