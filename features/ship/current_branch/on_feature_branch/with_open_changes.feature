Feature: git town-ship: errors if there are open changes

  As a developer trying to ship a branch with uncommitted changes
  I should see an error that my branch is in an unfinished state
  So that my users don't experience half-baked features.

  Background:
    Given my repo has a feature branch named "feature"
    And my workspace has an uncommitted file
    And I am on the "feature" branch
    When I run "git-town ship"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH | LOCATION |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
