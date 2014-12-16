Feature: git ship: shipping the current feature branch

  As a developer having finished a feature
  I want to be able to ship it safely in one easy step
  So that I can quickly move on to the next feature and remain productive.


  Scenario: local feature branch
    Given I am on a local feature branch
    And the following commit exists in my repository
      | LOCATION | FILE NAME    | FILE CONTENT    |
      | local    | feature_file | feature content |
    When I run `git ship -m 'feature done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I am on a feature branch
    And the following commit exists in my repository
      | LOCATION | FILE NAME    | FILE CONTENT    |
      | remote   | feature_file | feature content |
    When I run `git ship -m 'feature done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
