Feature: git town-ship: shipping a child branch

  Background:
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    And my repo has a feature branch named "feature-3" as a child of "feature-2"
    And the following commits exist in my repo
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, remote | feature 1 commit | feature_1_file | feature 1 content |
      | feature-2 | local, remote | feature 2 commit | feature_2_file | feature 2 content |
      | feature-3 | local, remote | feature 3 commit | feature_3_file | feature 3 content |
    And I am on the "feature-3" branch
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-3 | git fetch --prune --tags |
    And it prints the error:
      """
      shipping this branch would ship "feature-1, feature-2" as well,
      please ship "feature-1" first
      """
    And I am now on the "feature-3" branch
    And my repo is left with my original commits
    And Git Town is still aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |
      | feature-3 | feature-2 |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature-3" branch
    And my repo is left with my original commits
