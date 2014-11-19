Feature: Git Kill: Killing the current local feature branch

  Background:
    Given I have a feature branch named "good-feature"
    And I am on the local "unfortunate" branch
    And the following commits exist in my repository
      | branch       | location         | message            | file name        |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local            | unfortunate commit | unfortunate_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
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


  Scenario: Undoing a kill of a local feature branch
    When I run `git kill --undo`
    Then I end up on the "unfortunate" branch
    And the existing branches are
      | repository | branches                        |
      | local      | main, unfortunate, good-feature |
      | remote     | main, good-feature              |
    And I have the following commits
      | branch       | location         | message            | files            |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local            | unfortunate commit | unfortunate_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
