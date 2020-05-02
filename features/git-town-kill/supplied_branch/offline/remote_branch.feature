Feature: git town-kill: killing a remote branch in offline mode

  When offline and trying to kill a remote branch
  I want to be notified that this operation is not possible
  So that I know about my mistake and can do more appropriate actions instead.


  Background:
    Given Git Town is in offline mode
    And my origin has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | remote   | feature commit |
    And my repository knows about the remote branch
    And I am on the "main" branch
    When I run "git-town kill feature"


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND |
    And it prints the error:
      """
      Cannot delete remote branch 'feature' in offline mode
      """
    And I am still on the "main" branch
