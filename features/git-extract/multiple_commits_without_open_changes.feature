Feature: git extract: extracting multiple commits (with open changes)

  As a developer having a feature branch with many commits around an unrelated issue
  I want to be able to extract the commits into their own branch
  So that the issue can be reviewed separately and my feature branch remains focussed


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor1 commit   | refactor1_file   |
      |         |          | refactor2 commit   | refactor2_file   |
    When I run `git extract refactor` with the last two commit shas


  Scenario: result
    Then I end up on the "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILES            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
    And now I have the following committed files
      | BRANCH   | FILES                                            |
      | main     | remote_main_file                                 |
      | feature  | feature_file, refactor1_file, refactor2_file     |
      | refactor | remote_main_file, refactor1_file, refactor2_file |
