Feature: Aborting by leaving the commit message blank

  Scenario: Blank commit message
    Given I am on the "feature" branch
    And the following commit exists in my repository
      | location | message  | file name    | file content    |
      | local    | a commit | feature_file | feature content |
    And I run `git ship` with an empty commit message
    Then I am still on the "feature" branch
    And I still have the following commits
      | branch  | location | message  | files        |
      | feature | local    | a commit | feature_file |
    And I still have the following committed files
      | branch  | files        | content         |
      | feature | feature_file | feature content |

