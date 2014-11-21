Feature: Git Kill: Killing the current feature branch with open changes

  Background:
    Given I have feature branches named "good-feature" and "unfortunate"
    And the following commits exist in my repository
      | branch       | location         | message            | file name        |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
    And I am on the "unfortunate" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then I end up on the "unfortunate" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | repository | branches                        |
      | local      | main, unfortunate, good-feature |
      | remote     | main, unfortunate, good-feature |
    And I have the following commits
      | branch       | location         | message            | files            |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
