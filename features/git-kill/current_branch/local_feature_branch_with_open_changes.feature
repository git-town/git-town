Feature: Git Kill: Killing the current local feature branch

  Background:
    Given I have a feature branch named "good-feature"
    And I have a local feature branch named "dead-end-feature"
    And the following commits exist in my repository
      | branch           | location         | message         | file name        |
      | good-feature     | local and remote | good commit     | good_file        |
      | dead-end-feature | local            | dead-end commit | unfortunate_file |
    And I am on the "dead-end-feature" branch
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
    Then I end up on the "dead-end-feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | repository | branches                             |
      | local      | main, dead-end-feature, good-feature |
      | remote     | main, good-feature                   |
    And I have the following commits
      | branch           | location         | message         | files            |
      | good-feature     | local and remote | good commit     | good_file        |
      | dead-end-feature | local            | dead-end commit | unfortunate_file |
