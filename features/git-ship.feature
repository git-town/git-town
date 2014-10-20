Feature: Git Ship

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


  Scenario: feature branch not ahead of main
    Given I am on a feature branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The current branch 'feature' has no commits to merge into 'main'."
    And I end up on the "feature" branch


  Scenario: on the main branch
    Given I am on the main branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The current branch 'main' is not a feature branch. Please checkout a feature branch to ship"
    And I am still on the "main" branch
    And there are no commits
    And there are no open changes


  Scenario: on non feature branch
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The current branch 'production' is not a feature branch. Please checkout a feature branch to ship"
    And I am still on the "production" branch
    And there are no commits
    And there are no open changes


  Scenario: with uncommitted changes
    Given I am on a feature branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "You should not ship while having open files in Git"
    And I am still on the feature branch
    And there are no commits
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after conflict while pulling the feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | location | message                   | file name        | file content   |
      | remote   | conflicting remote commit | conflicting_file | remote content |
      | local    | conflicting local commit  | conflicting_file | local content  |
    When I run `git ship -m 'feature done'` while allowing errors
    Then my repo has a merge in progress
    And there is an abort script for "git ship"
    When I run `git ship --abort`
    Then I end up on my feature branch
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And there are no open changes
    And my branch and its remote still have 1 and 1 different commits


  Scenario: user aborts after conflict while pulling the main branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch | location | message                   | file name        | file content   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
    When I run `git ship -m 'feature done'` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git ship"
    When I run `git ship --abort`
    Then I end up on my feature branch
    And there is no rebase in progress
    And there is no abort script for "git ship" anymore
    And there are no open changes
    And the "main" branch and its remote still have 1 and 1 different commits


  Scenario: user aborts after conflict while merging the main branch into the feature
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                    | file name        | file content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
    When I run `git ship -m 'feature done'` while allowing errors
    Then I end up on the "feature" branch
    And file "conflicting_file" has a merge conflict
    And there is an abort script for "git ship"
    When I run `git ship --abort`
    Then I end up on the feature branch
    And there are no merge conflicts anymore
    And there is no abort script for "git ship" anymore
    And now I have the following commits
      | branch  | location | message                     | files            |
      | main    | local    | conflicting main commit     | conflicting_file |
      | feature | local    | conflicting feature commit  | conflicting_file |

