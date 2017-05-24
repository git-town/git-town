Feature: git town-ship: errors if on supplied branch and there are open changes

  (see ../../../current_branch/on_feature_branch/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And I have an uncommitted file
    And I am on the "feature" branch
    When I run `git-town ship feature`


  Scenario: result
    Then it runs no commands
    And I get the error "You have uncommitted changes. Did you mean to commit them before shipping?"
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And there are no commits
