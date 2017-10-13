Feature: git town-ship: shipping a child branch

  As a user shipping a feature branch that is a child branch of another unshipped feature branch
  I want to see a warning that this branch has unshipped parents
  So that I don't accidentally also ship the parent branch.


  Background:
    Given my repository has a feature branch named "feature-1"
    And my repository has a feature branch named "feature-2" as a child of "feature-1"
    And it has a feature branch named "feature-3" as a child of "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local and remote | feature 1 commit | feature_1_file | feature 1 content |
      | feature-2 | local and remote | feature 2 commit | feature_2_file | feature 2 content |
      | feature-3 | local and remote | feature 3 commit | feature_3_file | feature 3 content |
    And I am on the "feature-3" branch
    When I run `git-town ship`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND           |
      | feature-3 | git fetch --prune |
    And it prints the error "Shipping this branch would ship feature-1, feature-2 as well."
    And it prints the error "Please ship "feature-1" first."
    And I end up on the "feature-3" branch
    And my repository is left with my original commits
    And my branch hierarchy metadata is unchanged


  Scenario: undo
    When I run `git-town ship --undo`
		Then Git Town runs no commands
    And it prints the error "Nothing to undo"
    And I am still on the "feature-3" branch
    And my repository is left with my original commits
