Feature: Git Kill: killing the given feature branch without open changes

  Background:
    Given I have feature branches named "good-feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "good-feature" branch
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                       |
      | good-feature | git fetch --prune             |
      | good-feature | git push origin :dead-feature |
      | good-feature | git branch -D dead-feature    |
    And I am still on the "good-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH       | COMMAND                                       |
      | good-feature | git branch dead-feature [SHA:dead-end commit] |
      | good-feature | git push -u origin dead-feature               |
    And I am still on the "good-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                         |
      | local      | main, dead-feature, good-feature |
      | remote     | main, dead-feature, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILES            |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
