Feature: git ship: shipping a child branch

  As a user shipping a feature branch that is a child branch of another unshipped feature branch
  I want to see a warning that this branch has unshipped parents
  So that I don't accidentally also ship the parent branch.


  Background:
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
    And I am on the "child-feature" branch
    When I run `git ship -m "child feature done"`


  Scenario: result
    Then I get the error "Shipping this branch would ship 'parent-feature' as well."
    And I get the error "Please ship 'parent-feature' first."
    And it runs no Git commands
    And I end up on the "child-feature" branch
    And I am left with my original commits
    And my branch hierarchy metadata is unchanged

