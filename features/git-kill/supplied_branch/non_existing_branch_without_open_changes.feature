Feature: git kill: don't delete a non-existing branch (without open changes)

  (see ./non_existing_branch_with_open_changes.feature)


  Background:
    Given I am on the "good-feature" branch
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
    When I run `git kill non-existing-feature` while allowing errors

  Scenario: result
    Then I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "good-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |
