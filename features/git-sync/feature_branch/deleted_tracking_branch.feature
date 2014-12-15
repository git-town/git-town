Feature: git-sync: restore a deleted tracking branch

  As a developer on a feature branch whose remote branch was deleted by somebody else
  I want a new tracking branch to be created for my local branch
  So that my work is safe in case my local copy gets lost.


  Scenario: without a remote branch
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | feature | local and remote | feature commit | feature_file |
    And the "feature" branch gets deleted on the remote
    And I am on the "feature" branch
    When I run `git sync`
    Then I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE        | FILES        |
      | feature | local and remote | feature commit | feature_file |
