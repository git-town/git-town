Feature: Git Ship: shipping the supplied feature branch with open changes


  Scenario: local feature branch
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | main    | local    | main commit    | main_file    | main content    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "main_file" and content: "conflicting content"
    When I run `git ship feature -m 'feature done'`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "main_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | main commit  | main_file    |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES                   |
      | main   | feature_file, main_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship feature -m 'feature done'`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "feature_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
