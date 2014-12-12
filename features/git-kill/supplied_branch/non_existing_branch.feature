Feature: git kill: does not delete the given non-existing branch

  As a developer trying to remove a non-existing feature branch
  I want the tool to recognize this situation and abort cleanly
  So that I become aware about what I'm doing wrong and can do it better next time.


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
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
