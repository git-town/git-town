Feature: git kill: don't delete a misspelled branch (with open changes)

  Background:
    Given I am on the "good-feature" branch
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill non-existing-feature` while allowing errors

  Scenario: result
    Then I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "good-feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |
