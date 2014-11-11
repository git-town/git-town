Feature: shipping the current feature branch


  Scenario: local feature branch
    Given I am on a local feature branch
    And the following commit exists in my repository
      | location | file name    | file content    |
      | local    | feature_file | feature content |
    When I run `git ship -m 'feature done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | branch  | location         | message      | files        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | branch | files        |
      | main   | feature_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I am on a feature branch
    And the following commit exists in my repository
      | location | file name    | file content    |
      | remote   | feature_file | feature content |
    When I run `git ship -m 'feature done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | branch  | location         | message      | files        |
      | main    | local and remote | feature done | feature_file |
    And now I have the following committed files
      | branch | files        |
      | main   | feature_file |