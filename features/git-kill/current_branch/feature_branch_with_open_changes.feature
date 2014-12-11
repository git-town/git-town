Feature: Git Kill: Killing the current feature branch with open changes

  As a developer working on a dead-end feature branch
  I want to be able to cleanly delete the whole current branch including open changes
  So that my workspace is ready to work on something else and my productivity remains high.


  Background:
    Given I have feature branches named "good-feature" and "dead-end-feature"
    And the following commits exist in my repository
      | branch           | location         | message         | file name        |
      | good-feature     | local and remote | good commit     | good_file        |
      | dead-end-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "dead-end-feature" branch
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
    Then I end up on the "dead-end-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | repository | branches                        |
      | local      | main, dead-end-feature, good-feature |
      | remote     | main, dead-end-feature, good-feature |
    And I have the following commits
      | branch           | location         | message         | files            |
      | good-feature     | local and remote | good commit     | good_file        |
      | dead-end-feature | local and remote | dead-end commit | unfortunate_file |
