Feature: git ship: does'nt ship a feature branch without changes (without open changes)

  As a developer shipping a feature branch that doesn't result in any changes on main
  I want to be notified about this situation
  So that I can investigate and resolve it safely, don't ship the wrong things, and my users see meaningful progress.


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME   | FILE CONTENT   |
      | main    | remote   | common_file | common content |
      | feature | local    | common_file | common content |
    And I am on the "other_feature" branch
    When I run `git ship feature` while allowing errors


  Scenario: result
    Then I get the error "The branch 'feature' has no shippable changes"
    And I am still on the "other_feature" branch
