Feature: git town-ship: errors if there are open changes


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
