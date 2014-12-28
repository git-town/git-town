Feature: Git Kill: killing the current feature branch by name, without open changes

  Background:
    Given I have feature branches named "good-feature" and "delete-by-name"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE            | FILE NAME        |
      | good-feature   | local and remote | good commit        | good_file        |
      | delete-by-name | local and remote | unfortunate commit | unfortunate_file |
    And I am on the "delete-by-name" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill delete-by-name`


  Scenario: result
    Then I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then I end up on the "delete-by-name" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                           |
      | local      | main, delete-by-name, good-feature |
      | remote     | main, delete-by-name, good-feature |
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE            | FILES            |
      | good-feature   | local and remote | good commit        | good_file        |
      | delete-by-name | local and remote | unfortunate commit | unfortunate_file |
