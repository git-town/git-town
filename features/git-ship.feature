Feature: Git Ship

  Scenario: local feature branch
    Given I am on a local feature branch
    And the following commit exists
      | location | file name    | file content    |
      | local    | feature_file | feature content |
    When I run `git ship 'feature_done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | branch  | message      | files        |
      | main    | feature_done | feature_file |
    And now I have the following committed files
      | branch | name         | content         |
      | main   | feature_file | feature content |


  Scenario: feature branch with non-pulled updates in the repo
    Given I am on a feature branch
    And the following commit exists
      | location | file name    | file content    |
      | remote   | feature_file | feature content |
    When I run `git ship 'feature_done'`
    Then I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | branch  | message      | files        |
      | main    | feature_done | feature_file |
    And now I have the following committed files
      | branch | name         | content         |
      | main   | feature_file | feature content |


  Scenario: on the main branch
    Given I am on the main branch
    When I run `git ship 'feature_done'` while allowing errors
    Then I get the error "Please checkout the feature branch to ship"
    And I am still on the "main" branch
    And there are no commits
    And there are no open changes


  Scenario: with uncommitted changes
    Given I am on a feature branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship 'feature_done'` while allowing errors
    Then I get the error "You should not ship while having open files in Git"
    And I am still on the feature branch
    And there are no commits
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: conflict after pulling the feature branch
    Given I am on a feature branch
    And the following commits exist
      | location | message                   | file name        | file content   |
      | remote   | conflicting remote commit | conflicting_file | remote content |
      | local    | conflicting local commit  | conflicting_file | local content  |
    When I run `git ship 'feature_done'` while allowing errors
    Then I get the error "ERROR WHILE PULLING THE FEATURE BRANCH"
    And my repo has a rebase in progress
    And there is an abort script for "git ship"
    When I run `git ship --abort`
    Then I end up on my feature branch
    And there is no rebase in progress
    And there is no abort script for "git ship" anymore
    And there are no open changes
    And my branch and its remote still have 1 and 1 different commits


  Scenario: conflict after the squash-merge of the feature branch into the main branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                    | file name        | file content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
    When I run `git ship 'feature_done'` while allowing errors
    Then I get the error "ERROR WHILE SQUASH-MERGING THE FEATURE BRANCH"
    And I end up on the "main" branch
    And file "conflicting_file" has a merge conflict
    And there is an abort script for "git ship"
    When I run `git ship --abort`
    Then I end up on the feature branch
    And there are no merge conflicts anymore
    And there is no abort script for "git ship" anymore
    And now I have the following commits
      | branch  | message                     | files            |
      | main    | conflicting main commit     | conflicting_file |
      | feature | conflicting feature commit  | conflicting_file |

