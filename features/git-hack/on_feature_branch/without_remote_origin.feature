Feature: git hack: starting a new feature from a feature branch (without remote repo)

  As a developer working on something unrelated to my current feature branch and without a remote repository
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given I have a feature branch named "existing-feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH           | LOCATION | MESSAGE                 | FILE NAME    |
      | main             | local    | main commit             | main_file    |
      | existing-feature | local    | existing feature commit | feature_file |
    And I am on the "existing-feature" branch
    And I have an uncommitted file
    When I run `git hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                          |
      | existing-feature | git stash -u                     |
      |                  | git checkout -b new-feature main |
      | new-feature      | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH           | LOCATION | MESSAGE                 | FILE NAME    |
      | main             | local    | main commit             | main_file    |
      | existing-feature | local    | existing feature commit | feature_file |
      | new-feature      | local    | main commit             | main_file    |
