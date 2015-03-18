Feature: git hack: starting a new feature from a feature branch (without open changes or remote repo)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "existing_feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH           | LOCATION | MESSAGE                 | FILE NAME    |
      | main             | local    | main commit             | main_file    |
      | existing_feature | local    | existing feature commit | feature_file |
    And I am on the "existing_feature" branch
    When I run `git hack new_feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                          |
      | existing_feature | git checkout -b new_feature main |
    And I end up on the "new_feature" branch
    And I have the following commits
      | BRANCH           | LOCATION | MESSAGE                 | FILE NAME    |
      | main             | local    | main commit             | main_file    |
      | existing_feature |          | existing feature commit | feature_file |
      | new_feature      |          | main commit             | main_file    |
