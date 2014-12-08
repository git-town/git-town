Feature: Git Ship: errors when with uncommitted changes

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    |
      | feature | local    | feature commit | feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship -m 'feature done'` while allowing errors

  Scenario: result
    Then I get the error "You cannot ship with uncommitted changes."
    And I am still on the feature branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
