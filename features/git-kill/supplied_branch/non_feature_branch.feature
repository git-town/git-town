Feature: Git Kill: Non-feature branches are not killed


  Background:
    Given I have a feature branch named "good-feature"
    Given non-feature branch configuration "qa"
    And the following commits exist in my repository
      | branch       | location         | message     | file name |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
    And I am on the "good-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill qa` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "good-feature" branch
    And the existing branches are
      | repository | branches               |
      | local      | main, qa, good-feature |
      | remote     | main, qa, good-feature |
    And I have the following commits
      | branch | location         | message   | files   |
      | qa     | local and remote | qa commit | qa_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
