Feature: Git Kill: killing the given feature branch with open changes

  Background:
    Given I have a feature branch named "good-feature"
    And I have a feature branch named "delete-by-name-feature"
    And the following commits exist in my repository
      | branch                 | location         | message            | file name        |
      | good-feature           | local and remote | good commit        | good_file        |
      | delete-by-name-feature | local and remote | unfortunate commit | unfortunate_file |
    And I am on the "good-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill delete-by-name-feature`


  Scenario: result
    Then I am still on the "good-feature" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then I end up on the "delete-by-name-feature" branch
    And the existing branches are
      | repository | branches                                   |
      | local      | main, delete-by-name-feature, good-feature |
      | remote     | main, delete-by-name-feature, good-feature |
    And I have the following commits
      | branch                 | location         | message            | files            |
      | good-feature           | local and remote | good commit        | good_file        |
      | delete-by-name-feature | local and remote | unfortunate commit | unfortunate_file |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"

