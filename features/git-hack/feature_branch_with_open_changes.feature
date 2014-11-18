Feature: hacking on a feature branch with open changes


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | branch  | location | message        | file name    |
      | main    | remote   | main commit    | main_file    |
      | feature | local    | feature commit | feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack other_feature`


  Scenario: result
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | branch        | location         | message        | files        |
      | main          | local and remote | main commit    | main_file    |
      | feature       | local            | feature commit | feature_file |
      | other_feature | local            | main commit    | main_file    |
