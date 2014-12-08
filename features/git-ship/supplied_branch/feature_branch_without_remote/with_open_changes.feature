Feature: Git Ship: feature branch without a remote with open changes

  Background:
    Given I have local feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    |
      | feature | local    | feature commit | feature_file |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship feature -m 'feature done'`


  Scenario: result
    Then I am still on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the branch "feature" is deleted everywhere
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
