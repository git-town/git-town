Feature: git town-ship: errors if on supplied branch and there are open changes

  (see ../../../current_branch/on_feature_branch/with_open_changes.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And my workspace has an uncommitted file
    And I am on the "feature" branch
    When I run "git-town ship feature"


  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      You have uncommitted changes. Did you mean to commit them before shipping?
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repository now has the following commits
      | BRANCH | LOCATION |
