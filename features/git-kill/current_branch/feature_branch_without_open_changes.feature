Feature: Git Kill: Killing the current feature branch without open changes

  Background:
    Given I have a feature branch named "good-feature"
    And I am on the "unfortunate" branch
    And the following commits exist in my repository
      | branch       | location         | message            | file name        |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
    When I run `git kill`


  Scenario: result
    Then I end up on the "main" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: Undoing the kill
    When I run `git kill --undo`
    Then I end up on the "unfortunate" branch
    And the existing branches are
      | repository | branches                                |
      | local      | main, unfortunate, good-feature |
      | remote     | main, unfortunate, good-feature |
    And I have the following commits
      | branch       | location         | message            | files            |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
    And the branch "good-feature" still exists
