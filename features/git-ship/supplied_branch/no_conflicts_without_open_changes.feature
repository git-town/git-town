Feature: Git Ship: shipping the supplied feature branch without open changes


  Scenario: local feature branch
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature_file | feature content |
    And I am on the "other_feature" branch
    When I run `git ship feature -m 'feature done'`
    Then I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILES        |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "other_feature" branch
    When I run `git ship feature -m 'feature done'`
    Then I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILES        |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
