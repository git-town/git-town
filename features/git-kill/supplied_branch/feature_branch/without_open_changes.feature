Feature: git kill: killing the given feature branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have feature branches named "feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "feature" branch
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                       |
      | feature | git fetch --prune             |
      | feature | git push origin :dead-feature |
      | feature | git branch -D dead-feature    |
    And I am still on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                       |
      | feature | git branch dead-feature [SHA:dead-end commit] |
      | feature | git push -u origin dead-feature               |
    And I am still on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
      | remote     | main, dead-feature, feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
