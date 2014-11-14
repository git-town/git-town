Feature: Git Kill: The main branch is not killed


  Background:
    Given I have a feature branch named "good-feature"
    And I am on the "main" branch
    And the following commits exist in my repository
      | branch              | location         | message            | file name        |
      | good-feature        | local and remote | good commit        | good_file        |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

