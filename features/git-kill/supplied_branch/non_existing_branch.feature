Feature: Git Kill: does not delete the given non-existing branch

  Background:
    Given I am on the "good-feature" branch
    And the following commits exist in my repository
      | branch       | location         | message     | file name |
      | good-feature | local and remote | good commit | good_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill non-existing-feature` while allowing errors

  Scenario: result
    Then I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "good-feature" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

