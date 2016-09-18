Feature: git town-ship: shipping a child branch

  (see ../../current_branch/on_feature_branch/on_child_branch.feature)


  Background:
    Given I have a feature branch named "feature-1"
    And I have a feature branch named "feature-2" as a child of "feature-1"
    And I have a feature branch named "feature-3" as a child of "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local and remote | feature 1 commit | feature_1_file | feature 1 content |
      | feature-2 | local and remote | feature 2 commit | feature_2_file | feature 2 content |
      | feature-3 | local and remote | feature 3 commit | feature_3_file | feature 3 content |
    And I am on the "feature-1" branch
    When I run `git town-ship feature-3 -m "feature 3 done"`


  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND           |
      | feature-1 | git fetch --prune |
    And I get the error "Shipping this branch would ship "feature-1", "feature-2" as well."
    And I get the error "Please ship "feature-1" first."
    And I end up on the "feature-1" branch
    And I am left with my original commits
    And my branch hierarchy metadata is unchanged
