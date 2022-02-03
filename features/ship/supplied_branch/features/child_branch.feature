Feature: does not ship a child branch

  Background:
    Given my repo has a feature branch "feature-1"
    And my repo has a feature branch "feature-2" as a child of "feature-1"
    And my repo has a feature branch "feature-3" as a child of "feature-2"
    And my repo contains the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, remote | feature 1 commit |
      | feature-2 | local, remote | feature 2 commit |
      | feature-3 | local, remote | feature 3 commit |
    And I am on the "feature-1" branch
    When I run "git-town ship feature-3 -m 'feature 3 done'"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git fetch --prune --tags |
    And it prints the error:
      """
      shipping this branch would ship "feature-1, feature-2" as well,
      please ship "feature-1" first
      """
    And I am now on the "feature-1" branch
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature-1" branch
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
