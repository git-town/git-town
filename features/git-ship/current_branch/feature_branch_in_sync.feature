Feature: Git Ship: feature branch in sync with remote

  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | feature | local and remote | feature commit | feature_file |
    And I run `git ship -m 'feature done'`


  Scenario: result
    Then I end up on the "main" branch
    And there are no more feature branches
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
